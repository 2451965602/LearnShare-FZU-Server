package service

import (
    "LearnShare/biz/dal/db"
    model "LearnShare/biz/model/module"
    "LearnShare/biz/model/audit"
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

    reviews, err := db.GetPendingResourceReviews(s.ctx, int(req.PageNum), int(req.PageSize))
    if err != nil {
        return nil, err
    }

    var moduleReviews []*model.Review
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