package oss

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"LearnShare/config"
	"LearnShare/pkg/errno"
	"LearnShare/pkg/logger"

	"github.com/h2non/filetype"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
)

// IsFile 检查文件类型，allowedTypes 为空表示不做白名单限制
func IsFile(data *multipart.FileHeader, allowedTypes []string) error {
	file, err := data.Open()
	if err != nil {
		return errno.NewErrNo(errno.IOOperateErrorCode, "打开文件失败")
	}
	defer file.Close()

	// 读取文件头
	buffer := make([]byte, 512)
	n, err := io.ReadFull(file, buffer)
	if err != nil && err != io.EOF && !errors.Is(err, io.ErrUnexpectedEOF) {
		return errno.NewErrNo(errno.IOOperateErrorCode, "读取文件失败")
	}
	if n < 12 {
		return errno.NewErrNo(errno.ParamVerifyErrorCode, "文件体积过小")
	}

	kind, _ := filetype.Match(buffer[:n])

	if len(allowedTypes) == 0 {
		return nil
	}

	if kind != filetype.Unknown {
		for _, t := range allowedTypes {
			if kind.MIME.Value == t {
				return nil
			}
		}
		logger.Debugf("Upload type mismatch on content: detected=%s allowed=%v filename=%s", kind.MIME.Value, allowedTypes, data.Filename)
	}

	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(data.Filename)), ".")
	var extMime string
	switch ext {
	case "pdf":
		extMime = "application/pdf"
	case "docx":
		extMime = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case "pptx":
		extMime = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	case "zip":
		extMime = "application/zip"
	}
	if extMime != "" {
		for _, t := range allowedTypes {
			if t == extMime {
				return nil
			}
		}
	}
	logger.Debugf("Upload type rejected: detected=%s ext=%s extMime=%s allowed=%v filename=%s", func() string {
		if kind == filetype.Unknown {
			return "unknown"
		} else {
			return kind.MIME.Value
		}
	}(), ext, extMime, allowedTypes, data.Filename)
	return errno.NewErrNo(errno.ParamVerifyErrorCode, "不支持的文件类型")
}

// SaveFile 保存上传文件到本地临时目录，返回本地完整路径
func SaveFile(data *multipart.FileHeader, storePath, fileName string) (string, error) {
	// 确保目录存在
	if err := os.MkdirAll(storePath, 0o755); err != nil {
		return "", errno.NewErrNo(errno.IOOperateErrorCode, "创建目录失败")
	}

	// 确保文件名安全
	fileName = filepath.Base(fileName)
	fullPath := filepath.Join(storePath, fileName)

	// 打开目标文件
	dist, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return "", errno.NewErrNo(errno.IOOperateErrorCode, "打开文件失败")
	}
	defer func() { _ = dist.Close() }()

	// 打开上传文件
	src, err := data.Open()
	if err != nil {
		return "", errno.NewErrNo(errno.IOOperateErrorCode, "读取上传文件失败")
	}
	defer func() { _ = src.Close() }()

	// 复制并同步
	if _, err = io.Copy(dist, src); err != nil {
		return "", errno.NewErrNo(errno.IOOperateErrorCode, "写入文件失败")
	}
	if err = dist.Sync(); err != nil {
		// 同步失败不一定要当作致命错误，但记录并返回
		return "", errno.NewErrNo(errno.IOOperateErrorCode, "文件同步失败")
	}
	return fullPath, nil
}

// Upload 上传本地文件到七牛 OSS，返回外链
func Upload(localFile, filename, class string, targetId int64) (string, error) {
	key := fmt.Sprintf("%v/%v/%v", class, targetId, filename)

	putPolicy := storage.PutPolicy{
		Scope:      fmt.Sprintf("%s:%s", config.Oss.BucketName, key),
		InsertOnly: 0,
	}

	mac := auth.New(config.Oss.AccessKeyID, config.Oss.AccessKeySecret)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	// 自动探测区域：当配置为空、为 auto 或者配置不匹配时由 SDK 动态解析
	if strings.EqualFold(config.Oss.Zone, "") || strings.EqualFold(config.Oss.Zone, "auto") {
		if region, err := storage.GetRegion(config.Oss.AccessKeyID, config.Oss.BucketName); err == nil {
			cfg.Region = region
		} else {
			cfg.Region = getQiniuZone(config.Oss.Zone)
		}
	} else {
		cfg.Region = getQiniuZone(config.Oss.Zone)
	}
	cfg.UseCdnDomains = config.Oss.UseCdnDomains // 由配置决定

	resumeUploader := storage.NewResumeUploaderV2(&cfg)
	ret := storage.PutRet{}

	recorder, err := storage.NewFileRecorder(os.TempDir())
	if err != nil {
		return "", errno.NewErrNo(errno.QiNiuYunFileErrorCode, "创建断点记录器失败")
	}

	putExtra := storage.RputV2Extra{
		Recorder: recorder,
	}

	// 上传操作使用带超时的 Context
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if err = resumeUploader.PutFile(ctx, &ret, upToken, key, localFile, &putExtra); err != nil {
		// 尝试删除临时文件（即使上传失败也应清理）
		_ = os.Remove(localFile)
		return "", errno.NewErrNo(errno.QiNiuYunFileErrorCode, fmt.Sprintf("上传失败: %v", err))
	}

	// 上传成功后删除本地临时文件
	if rmErr := os.Remove(localFile); rmErr != nil {
		logger.Error(errno.NewErrNo(errno.QiNiuYunFileErrorCode, "删除本地临时文件失败"))
	}

	return storage.MakePublicURL(config.Oss.Endpoint, ret.Key), nil
}

func getQiniuZone(region string) *storage.Region {
	switch region {
	case "z0":
		return &storage.Zone_z0
	case "z1":
		return &storage.Zone_z1
	case "z2":
		return &storage.Zone_z2
	case "na0":
		return &storage.Zone_na0
	case "as0":
		return &storage.Zone_as0
	default:
		return &storage.Zone_z0
	}
}

// UploadFile 接收 multipart 文件并完成校验、存储和上传
func UploadFile(data *multipart.FileHeader, class string, targetId int64) (string, error) {
	var allowedTypes []string
	if class == "avatar" {
		allowedTypes = []string{"image/jpeg", "image/png", "image/gif", "image/webp"}
	} else if class == "resource" {
		allowedTypes = []string{
			"application/pdf",
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			"application/vnd.openxmlformats-officedocument.presentationml.presentation",
			"application/zip",
		}
		// 资源文件：若扩展名在白名单中，直接放行，避免部分环境内容类型识别不稳定
		ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(data.Filename)), ".")
		switch ext {
		case "pdf", "docx", "pptx", "zip":
			// 直接通过，不再依赖内容识别
		default:
			// 扩展名不在白名单，退回到内容识别
			if err := IsFile(data, allowedTypes); err != nil {
				return "", errno.NewErrNo(errno.ParamVerifyErrorCode, "不支持的文件类型")
			}
		}
	} else {
		// 其他分类仍使用内容识别
		if err := IsFile(data, allowedTypes); err != nil {
			return "", errno.NewErrNo(errno.ParamVerifyErrorCode, "不支持的文件类型")
		}
	}

	if class == "resource" {
		if data.Size > 3*1024*1024 {
			return "", errno.NewErrNo(errno.ParamVerifyErrorCode, "文件大小超过限制")
		}
	}

	// 使用系统临时目录，避免写入项目目录
	storePath := filepath.Join(os.TempDir(), class, strconv.FormatInt(targetId, 10))
	fileName := fmt.Sprintf("%v_%v_%v", targetId, generateRandomString(10), filepath.Base(data.Filename))

	localPath, err := SaveFile(data, storePath, fileName)
	if err != nil {
		return "", err
	}

	link, err := Upload(localPath, fileName, class, targetId)
	if err != nil {
		return "", err
	}
	return link, nil
}

// DeleteByURL 根据文件外链删除七牛云上的文件
func DeleteByURL(fileURL string) error {
	if fileURL == "" {
		return errno.NewErrNo(errno.ParamVerifyErrorCode, "空的 URL")
	}

	u, err := url.Parse(fileURL)
	if err != nil {
		return errno.NewErrNo(errno.ParamVerifyErrorCode, "无效的 URL")
	}

	key := strings.TrimPrefix(u.Path, "/")
	if key == "" {
		return errno.NewErrNo(errno.ParamVerifyErrorCode, "无法解析出文件 key")
	}

	mac := auth.New(config.Oss.AccessKeyID, config.Oss.AccessKeySecret)
	cfg := storage.Config{}
	cfg.Region = getQiniuZone(config.Oss.Zone)

	bm := storage.NewBucketManager(mac, &cfg)
	if err := bm.Delete(config.Oss.BucketName, key); err != nil {
		return errno.NewErrNo(errno.QiNiuYunFileErrorCode, fmt.Sprintf("删除失败: %v", err))
	}
	return nil
}

// generateRandomString 使用 crypto/rand 生成安全的 base62 字符串
func generateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	if n <= 0 {
		return ""
	}

	// 生成随机字节
	out := make([]byte, n)
	// 为高效使用少量随机数据，采用直接填充方式
	buf := make([]byte, 8)
	for i := 0; i < n; i++ {
		if _, err := rand.Read(buf); err != nil {
			// 退回到时间/大小相关的种子（极端情况）
			v := time.Now().UnixNano() + int64(i)
			out[i] = letters[int(math.Abs(float64(v)))%len(letters)]
			continue
		}
		v := binary.LittleEndian.Uint64(buf)
		out[i] = letters[int(v%uint64(len(letters)))]
	}
	return string(out)
}
