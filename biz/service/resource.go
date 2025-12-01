package service

import (
	"LearnShare/biz/dal/db"
	model "LearnShare/biz/model/module"
	"LearnShare/biz/model/resource"
	"LearnShare/pkg/errno"
	"LearnShare/pkg/logger"
	"LearnShare/pkg/oss"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"context"
	"mime/multipart"

	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"
)

// ResourceService 封装了资源相关的服务
type ResourceService struct {
	ctx context.Context
	c   *app.RequestContext
}

// NewResourceService 创建一个新的 ResourceService
func NewResourceService(ctx context.Context, c *app.RequestContext) *ResourceService {
	return &ResourceService{ctx: ctx, c: c}
}

// SearchResources 执行资源搜索
func (s *ResourceService) SearchResources(req *resource.SearchResourceReq) ([]*model.Resource, int64, error) {
	// 验证搜索关键词长度
	if req.Keyword != nil && *req.Keyword != "" && len(*req.Keyword) > 100 {
		return nil, 0, errno.ValidationKeywordTooLongError
	}

	// 验证分页参数
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}

	resources, total, err := db.SearchResources(s.ctx, req.Keyword, req.TagId, req.CourseID, req.SortBy, int(req.PageNum), int(req.PageSize))
	if err != nil {
		return nil, 0, err
	}

	modelResources := make([]*model.Resource, 0, len(resources))
	for _, r := range resources {
		modelResources = append(modelResources, r.ToResourceModule())
	}

	return modelResources, total, nil
}

// GetResource 执行获取单个资源信息
func (s *ResourceService) GetResource(req *resource.GetResourceReq) (*model.Resource, error) {
	// 验证资源ID
	if req.ResourceID <= 0 {
		return nil, errno.NewErrNo(errno.ServiceInvalidParameter, "资源ID无效")
	}

	resourcedata, err := db.GetResourceByID(s.ctx, req.ResourceID)
	if err != nil {
		return nil, err
	}

	return resourcedata.ToResourceModule(), nil
}

func (s *ResourceService) DownloadResource(req *resource.DownloadResourceReq) (string, error) {
	if req.ResourceID <= 0 {
		return "", errno.NewErrNo(errno.ServiceInvalidParameter, "资源ID无效")
	}

	r, err := db.GetResourceByID(s.ctx, req.ResourceID)
	if err != nil {
		return "", err
	}

	userID := GetUidFormContext(s.c)

	errChan := db.UpdateResourceAsync(s.ctx, req.ResourceID, map[string]interface{}{
		"download_count": gorm.Expr("download_count + 1"),
	})
	if err = <-errChan; err != nil {
		return "", err
	}

	repErrChan := db.IncrementUserReputationAsync(s.ctx, r.UploaderID, 1)
	if e := <-repErrChan; e != nil {
		return "", e
	}

	logger.Infof("user %d downloaded resource %d", userID, req.ResourceID)

	return r.FilePath, nil
}

// GetResourceComments 执行获取资源评论列表
func (s *ResourceService) GetResourceComments(req *resource.GetResourceCommentsReq) ([]*model.ResourceCommentWithUser, int64, error) {
	// 验证资源ID
	if req.ResourceID <= 0 {
		return nil, 0, errno.NewErrNo(errno.ServiceInvalidParameter, "资源ID无效")
	}

	// 验证分页参数
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}

	// 调用数据库层获取评论数据
	comments, total, err := db.GetResourceComments(s.ctx, req.ResourceID, req.SortBy, int(req.PageNum), int(req.PageSize))
	if err != nil {
		return nil, 0, err
	}

    modelComments := make([]*model.ResourceCommentWithUser, 0, len(comments))
    for _, comment := range comments {
        modelComments = append(modelComments, comment.ToResourceCommentWithUserModule())
    }

    return modelComments, total, nil
}

// SubmitResourceRating 执行提交资源评分
func (s *ResourceService) SubmitResourceRating(req *resource.SubmitResourceRatingReq) (*model.ResourceRating, error) {
	userID := GetUidFormContext(s.c)

	// 验证评分范围
	if req.Rating < 0 || req.Rating > 5 {
		return nil, errno.ValidationRatingRangeInvalidError
	}

	// 调用数据库层提交评分，使用rating字段
	rating, err := db.SubmitResourceRating(s.ctx, userID, req.ResourceID, req.Rating)
	if err != nil {
		return nil, err
	}

	return rating.ToResourceRatingModule(), nil
}

// SubmitResourceComment 执行提交资源评论
func (s *ResourceService) SubmitResourceComment(req *resource.SubmitResourceCommentReq) (*model.ResourceComment, error) {
	userID := GetUidFormContext(s.c)

	// 验证评论内容
	if req.Content == "" {
		return nil, errno.ResourceInvalidCommentError
	}

	if len(req.Content) > 1000 {
		return nil, errno.ValidationCommentTooLongError
	}

	// 处理父评论ID
	var parentID *int64
	if req.IsSetParentId() && req.ParentId != nil && *req.ParentId != 0 {
		parentID = req.ParentId
	}

	// 调用数据库层提交评论
	resultChan := db.SubmitResourceCommentAsync(s.ctx, userID, req.ResourceID, req.Content, parentID)
	result := <-resultChan
	if result.Err != nil {
		return nil, result.Err
	}

	if result.Comment == nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "提交资源评论失败: 未返回评论数据")
	}

	return result.Comment.ToResourceCommentModule(), nil
}

// DeleteResourceRating 执行删除资源评分
func (s *ResourceService) DeleteResourceRating(req *resource.DeleteResourceRatingReq) error {
	userID := GetUidFormContext(s.c)

	// 验证评分ID
	if req.RatingID <= 0 {
		return errno.NewErrNo(errno.ServiceInvalidParameter, "评分ID无效")
	}

	// 使用异步删除评分
	errChan := db.DeleteResourceRatingAsync(s.ctx, req.RatingID, userID)
	if err := <-errChan; err != nil {
		return err
	}

	return nil
}

// DeleteResourceComment 删除资源评论
func (s *ResourceService) DeleteResourceComment(req *resource.DeleteResourceCommentReq) error {
	userID := GetUidFormContext(s.c)

	// 验证评论ID
	if req.CommentID <= 0 {
		return errno.NewErrNo(errno.ServiceInvalidParameter, "评论ID无效")
	}

	// 使用异步删除评论
	errChan := db.DeleteResourceCommentAsync(s.ctx, req.CommentID, userID)
	if err := <-errChan; err != nil {
		return err
	}

	return nil
}

// ReportResource 举报一个资源
func (s *ResourceService) ReportResource(req *resource.ReportResourceReq) error {
	// 验证资源ID
	if req.ResourceID <= 0 {
		return errno.NewErrNo(errno.ServiceInvalidParameter, "资源ID无效")
	}

	// 验证举报原因
	if req.Reason == "" {
		return errno.ResourceReportInvalidReasonError
	}
	if len(req.Reason) > 500 {
		return errno.ValidationReportReasonTooLongError
	}

	// 从上下文获取当前用户ID
	userID := GetUidFormContext(s.c)

	// 使用异步创建举报记录
	errChan := db.CreateReviewAsync(s.ctx, userID, req.ResourceID, "resource", req.Reason)
	if err := <-errChan; err != nil {
		return err
	}

	return nil
}

// AdminDeleteResource 管理员硬删除资源
func (s *ResourceService) AdminDeleteResource(req *resource.AdminDeleteResourceReq) error {
	if req.ResourceID <= 0 {
		return errno.NewErrNo(errno.ServiceInvalidParameter, "资源ID无效")
	}

	if err := db.AdminDeleteResource(s.ctx, req.ResourceID); err != nil {
		return err
	}
	return nil
}

func (s *ResourceService) AdminDeleteResourceComment(req *resource.AdminDeleteResourceCommentReq) error {
	if req.CommentID <= 0 {
		return errno.NewErrNo(errno.ServiceInvalidParameter, "评论ID无效")
	}
	if err := db.AdminDeleteResourceComment(s.ctx, req.CommentID); err != nil {
		return err
	}
	return nil
}

func (s *ResourceService) AdminDeleteResourceRating(req *resource.AdminDeleteResourceRatingReq) error {
	if req.RatingID <= 0 {
		return errno.NewErrNo(errno.ServiceInvalidParameter, "评分ID无效")
	}
	if err := db.AdminDeleteResourceRating(s.ctx, req.RatingID); err != nil {
		return err
	}
	return nil
}

func (s *ResourceService) UploadResource(file *multipart.FileHeader, title string, description *string, courseID int64, tags []string) (*model.Resource, error) {
	if title == "" {
		return nil, errno.ParamVerifyError
	}
	if courseID <= 0 {
		return nil, errno.ParamVerifyError
	}
	if utf8.RuneCountInString(title) > 255 {
		return nil, errno.ParamVerifyError
	}
	if description != nil {
		if utf8.RuneCountInString(*description) > 500 {
			return nil, errno.ParamVerifyError
		}
	}
	for _, t := range tags {
		if strings.TrimSpace(t) == "" {
			continue
		}
		if utf8.RuneCountInString(strings.TrimSpace(t)) > 50 {
			return nil, errno.ParamVerifyError
		}
	}

	userID := GetUidFormContext(s.c)

	link, err := oss.UploadFile(file, "resource", courseID)
	if err != nil {
		return nil, err
	}

	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(file.Filename)), ".")
	switch ext {
	case "pdf", "docx", "pptx", "zip":
	default:
		return nil, errno.ParamVerifyError
	}

	res := &db.Resource{
		ResourceName: title,
		Description: func() string {
			if description != nil {
				return *description
			}
			return ""
		}(),
		FilePath:   link,
		FileType:   ext,
		FileSize:   file.Size,
		UploaderID: userID,
		CourseID:   courseID,
		Status:     "normal",
	}

	errChan := db.CreateResourceAsync(s.ctx, res)
	if err = <-errChan; err != nil {
		return nil, err
	}

	if len(tags) > 0 {
		// 使用批量操作优化标签处理
		tagsData, e := db.GetOrCreateTagsBatch(s.ctx, tags)
		if e != nil {
			return nil, e
		}

		// 提取所有 tagID
		tagIDs := make([]int64, len(tags))
		for i, tag := range tagsData {
			tagIDs[i] = tag.TagID
		}

		// 批量关联标签
		if e = db.LinkResourceTagsBatch(s.ctx, res.ResourceID, tagIDs); e != nil {
			return nil, e
		}
	}

	// 直接构建返回结果，避免重复查询
    var tagsResp []*model.ResourceTag
    if res.Tags != nil {
        for _, t := range res.Tags {
            tagsResp = append(tagsResp, t.ToResourceTagModule())
        }
    } else {
        tagsResp = make([]*model.ResourceTag, 0)
    }

    resp := res.ToResourceModule()
    resp.Tags = tagsResp

	return resp, nil
}

func (s *ResourceService) ReactResourceComment(commentID int64, action string) error {
	if commentID <= 0 {
		return errno.ParamVerifyError
	}
	switch action {
	case "like", "dislike", "cancel_like", "cancel_dislike":
	default:
		return errno.ParamVerifyError.WithMessage("操作类型无效")
	}
	userID := GetUidFormContext(s.c)
	errChan := db.ReactResourceCommentAsync(s.ctx, userID, commentID, action)
	if err := <-errChan; err != nil {
		return err
	}
	return nil
}
