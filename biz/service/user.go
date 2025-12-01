package service

import (
	"LearnShare/biz/dal/db"
	"LearnShare/biz/dal/redis"
	"LearnShare/biz/model/module"
	"LearnShare/biz/model/user"
	"LearnShare/pkg/errno"
	"LearnShare/pkg/logger"
	"LearnShare/pkg/oss"
	"LearnShare/pkg/utils"
	"context"
	"mime/multipart"
	"strconv"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"go.uber.org/zap"
)

type UserService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewUserService(ctx context.Context, c *app.RequestContext) *UserService {
	return &UserService{ctx: ctx, c: c}
}

func (s *UserService) Register(req *user.RegisterReq) error {

	if valid, err := utils.VerifyUsername(req.Username); !valid {
		logger.WithFields(
			zap.String("username", req.Username),
			zap.Error(err),
		).Warn("用户名验证失败")
		return err
	}

	if valid, err := utils.VerifyPassword(req.Password); !valid {
		return err
	}

	if valid, err := utils.VerifyEmail(req.Email); !valid {
		logger.WithFields(
			zap.String("email", logger.MaskEmail(req.Email)),
			zap.Error(err),
		).Warn("邮箱验证失败")
		return err
	}

	passwordHash, err := utils.EncryptPassword(req.Password)
	if err != nil {
		logger.WithFields(
			zap.String("username", req.Username),
			zap.Error(err),
		).Error("密码加密失败")
		return err
	}

	req.Password = passwordHash

	err = db.CreateUser(s.ctx, req.Username, req.Password, req.Email)
	if err != nil {
		logger.WithFields(
			zap.String("username", req.Username),
			zap.String("email", logger.MaskEmail(req.Email)),
			zap.Error(err),
		).Error("创建用户失败")
		return err
	}

	return nil
}

func (s *UserService) LoginIn(req *user.LoginInReq) (*module.User, error) {
	userInfo, err := db.GetUserByEmail(s.ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if err := utils.ComparePassword(userInfo.PasswordHash, req.Password); err != nil {
		logger.WithFields(
			zap.String("email", logger.MaskEmail(req.Email)),
			zap.Int64("user_id", userInfo.UserID),
		).Warn("登录失败：密码错误")
		return nil, err
	}

	userInfo.PasswordHash = ""
	return userInfo.ToUserModule(), nil
}

func (s *UserService) LoginOut() error {
	userId := GetUidFormContext(s.c)

	var errors []error

	uuidStr := GetUuidFormContext(s.c)
	if err := redis.SetBlacklistToken(s.ctx, uuidStr); err != nil {
		logger.WithFields(
			zap.Int64("user_id", userId),
			zap.Error(err),
		).Error("设置黑名单 token 失败")
		errors = append(errors, err)
	}

	// 如果有错误发生，记录日志但不阻止登出
	if len(errors) > 0 {
		return errors[0] // 返回第一个错误
	}

	return nil
}

func (s *UserService) SendVerifyEmail(req *user.SendVerifyEmailReq) error {
	code, err := utils.GenerateCode()
	if err != nil {
		return err
	}

	err = utils.MailSendCode(req.Email, code)
	if err != nil {
		return err
	}

	err = redis.PutCodeToCache(s.ctx, req.Email, code)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) VerifyEmail(req *user.VerifyEmailReq) error {

	userId := GetUidFormContext(s.c)

	if !redis.IsKeyExist(s.ctx, req.Email) {
		return errno.UserVerificationCodeExpiredError
	}

	storeCode, err := redis.GetCodeCache(s.ctx, req.Email)
	if err != nil {
		return err
	}

	if storeCode != req.Code {
		return errno.UserVerificationCodeInvalidError
	}

	// 使用异步更新用户状态
	errChan := db.UpdateUserStatuesAsync(s.ctx, userId, "active")
	if err := <-errChan; err != nil {
		return err
	}
	return nil
}

func (s *UserService) UpdateEmail(req *user.UpdateEmailReq) error {

	userId := GetUidFormContext(s.c)

	storeCode, err := redis.GetCodeCache(s.ctx, req.NewEmail)
	if err != nil {
		return err
	}

	if storeCode != req.Code {
		return errno.UserVerificationCodeInvalidError
	}

	// 使用异步更新用户邮箱
	errChan := db.UpdateUserEmailAsync(s.ctx, int(userId), req.NewEmail)
	if err := <-errChan; err != nil {
		return err
	}
	return nil
}

func (s *UserService) UpdatePassword(req *user.UpdatePasswordReq) error {
	userId := GetUidFormContext(s.c)
	userInfo, err := db.GetUserByID(s.ctx, userId)
	if err != nil {
		logger.WithFields(
			zap.Int64("user_id", userId),
			zap.Error(err),
		).Error("获取用户信息失败")
		return err
	}

	if err := utils.ComparePassword(userInfo.PasswordHash, req.OldPassword); err != nil {
		logger.WithFields(
			zap.Int64("user_id", userId),
		).Warn("修改密码失败：旧密码错误")
		return errno.UserPasswordIncorrectError
	}
	newPasswordHash, err := utils.EncryptPassword(req.NewPassword)
	if err != nil {
		logger.WithFields(
			zap.Int64("user_id", userId),
			zap.Error(err),
		).Error("新密码加密失败")
		return err
	}

	// 使用异步更新用户密码
	errChan := db.UpdateUserPasswordAsync(s.ctx, userInfo.UserID, newPasswordHash)
	if err := <-errChan; err != nil {
		logger.WithFields(
			zap.Int64("user_id", userId),
			zap.Error(err),
		).Error("更新用户密码失败")
		return err
	}

	return nil
}

func (s *UserService) UpdateMajor(req *user.UpdateMajorReq) error {
	userId := GetUidFormContext(s.c)
	var (
		userInfo *db.User
		err      error
	)

	if redis.IsKeyExist(s.ctx, strconv.FormatInt(userId, 10)) {
		userInfo, err = redis.GetUserInfoCache(s.ctx, strconv.FormatInt(userId, 10))
		if err != nil {
			return err
		}
	} else {
		userInfo, err = db.GetUserByID(s.ctx, userId)
		if err != nil {
			return err
		}
		err = redis.SetUserInfoCache(s.ctx, strconv.FormatInt(userId, 10), userInfo, 12*time.Hour)
		if err != nil {
			return err
		}
	}

	// 使用异步更新用户专业
	errChan := db.UpdateMajorIDAsync(s.ctx, int(userInfo.UserID), int(req.NewMajorId))
	if err := <-errChan; err != nil {
		return err
	}
	return nil
}

func (s *UserService) UploadAvatar(data *multipart.FileHeader) error {
	userId := GetUidFormContext(s.c)

	url, err := oss.UploadFile(data, "avatar", userId)
	if err != nil {
		return err
	}

	// 使用异步更新用户头像
	errChan := db.UpdateAvatarURLAsync(s.ctx, userId, url)
	if err = <-errChan; err != nil {
		return err
	}
	return nil
}

func (s *UserService) ResetPassword(req *user.ResetPasswordReq) error {
	storeCode, err := redis.GetCodeCache(s.ctx, req.Email)
	if err != nil {
		return err
	}

	if storeCode != req.Code {
		logger.WithFields(
			zap.String("email", logger.MaskEmail(req.Email)),
		).Warn("重置密码失败：验证码错误")
		return errno.UserVerificationCodeInvalidError
	}

	newPasswordHash, err := utils.EncryptPassword(req.NewPassword)
	if err != nil {
		logger.WithFields(
			zap.String("email", logger.MaskEmail(req.Email)),
			zap.Error(err),
		).Error("新密码加密失败")
		return err
	}

	userInfo, err := db.GetUserByEmail(s.ctx, req.Email)
	if err != nil {
		logger.WithFields(
			zap.String("email", logger.MaskEmail(req.Email)),
			zap.Error(err),
		).Error("获取用户信息失败")
		return err
	}

	// 使用异步更新用户密码
	errChan := db.UpdateUserPasswordAsync(s.ctx, userInfo.UserID, newPasswordHash)
	if err := <-errChan; err != nil {
		logger.WithFields(
			zap.String("email", logger.MaskEmail(req.Email)),
			zap.Int64("user_id", userInfo.UserID),
			zap.Error(err),
		).Error("重置密码失败")
		return err
	}

	return nil
}

func (s *UserService) GetUserInfo(req *user.GetUserInfoReq) (*module.User, error) {
	var (
		userInfo *db.User
		err      error
	)

	key := strconv.FormatInt(req.UserID, 10)

	if redis.IsKeyExist(s.ctx, key) {
		userInfo, err = redis.GetUserInfoCache(s.ctx, key)
		if err != nil {
			return nil, err
		}
		if userInfo != nil {
			return userInfo.ToUserModule(), nil
		}
		// 缓存命中但返回 nil，记录并回退到 DB
		logger.WithFields(zap.String("key", key)).Warn("缓存命中但返回 nil 用户信息，回退到数据库查询")
	}

	userInfo, err = db.GetUserByID(s.ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if userInfo == nil {
		return nil, errno.NewErrNo(errno.ServiceUserNotExist, "用户不存在")
	}

	if err = redis.SetUserInfoCache(s.ctx, key, userInfo, 12*time.Hour); err != nil {
		logger.WithFields(
			zap.Int64("user_id", req.UserID),
			zap.Error(err),
		).Warn("设置用户缓存失败")
		// 不阻塞，仍返回用户信息
	}

	return userInfo.ToUserModule(), nil
}
