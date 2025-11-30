package service

import (
	"LearnShare/biz/dal/db"
	"LearnShare/biz/model/audit"
	model "LearnShare/biz/model/module"
	"LearnShare/pkg/errno"

	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

// AuditService 封装了审核相关的服务
type AuditService struct {
	ctx context.Context
	c   *app.RequestContext
}

// NewAuditService 创建一个新的 AuditService
func NewAuditService(ctx context.Context, c *app.RequestContext) *AuditService {
	return &AuditService{ctx: ctx, c: c}
}

// GetResourceAuditList 获取待审核的资源举报列表
func (s *AuditService) GetResourceAuditList(req *audit.GetResourceAuditListReq) ([]*model.Review, error) {
	// 验证分页参数
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}

	// 统计总数用于越界页处理
	total, err := db.CountPendingResourceReviews(s.ctx)
	if err != nil {
		return nil, err
	}

	// 若请求页越界，自动定位到最后一页
	if total > 0 {
		maxPage := (total + int64(req.PageSize) - 1) / int64(req.PageSize)
		if int64(req.PageNum) > maxPage {
			req.PageNum = int32(maxPage)
		}
	}

	reviews, err := db.GetPendingResourceReviews(s.ctx, int(req.PageNum), int(req.PageSize))
	if err != nil {
		return nil, err
	}

	moduleReviews := make([]*model.Review, 0, len(reviews))
	for _, r := range reviews {
		moduleReviews = append(moduleReviews, r.ToReviewModule())
	}

	return moduleReviews, nil
}

// AuditResource 审核资源举报（approve/reject）
func (s *AuditService) AuditResource(req *audit.AuditResourceReq) error {
	// 验证参数
	if req.ReviewID <= 0 {
		return errno.NewErrNo(errno.ServiceInvalidParameter, "审核记录ID无效")
	}
	if req.Action != "approve" && req.Action != "reject" {
		return errno.NewErrNo(errno.ServiceInvalidParameter, "操作类型无效")
	}

	reviewerID := GetUidFormContext(s.c)
	// 调用数据库层执行审核
	err := db.AuditResourceReview(s.ctx, req.ReviewID, reviewerID, req.Action)
	if err != nil {
		return err
	}
	return nil
}

// GetResourceCommentAuditList 获取待审核的资源评论列表
func (s *AuditService) GetResourceCommentAuditList(req *audit.GetResourceCommentAuditListReq) ([]*model.Review, error) {
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}

	reviews, err := db.GetPendingCommentReviews(s.ctx, int(req.PageNum), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	reviewModules := make([]*model.Review, 0, len(reviews))
	for _, r := range reviews {
		reviewModules = append(reviewModules, r.ToReviewModule())
	}
	return reviewModules, nil
}

func (s *AuditService) AuditCourseComment(req *audit.AuditCourseCommentReq) error {
	if req.ReviewID <= 0 {
		return errno.NewErrNo(errno.ServiceInvalidParameter, "审核记录ID无效")
	}
	if req.Action != "approve" && req.Action != "reject" {
		return errno.NewErrNo(errno.ServiceInvalidParameter, "操作类型无效")
	}
	return nil
}

func (s *AuditService) AuditResourceComment(req *audit.AuditResourceCommentReq) error {
	if req.ReviewID <= 0 {
		return errno.NewErrNo(errno.ServiceInvalidParameter, "审核记录ID无效")
	}
	if req.Action != "approve" && req.Action != "reject" {
		return errno.NewErrNo(errno.ServiceInvalidParameter, "操作类型无效")
	}
	reviewerID := GetUidFormContext(s.c)
	if err := db.AuditResourceCommentReview(s.ctx, req.ReviewID, reviewerID, req.Action); err != nil {
		return err
	}
	return nil
}
