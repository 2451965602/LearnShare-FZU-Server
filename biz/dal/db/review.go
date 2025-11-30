package db

import (
	"LearnShare/pkg/constants"
	"LearnShare/pkg/errno"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

// GetPendingResourceReviews 获取待审核的资源举报列表
func GetPendingResourceReviews(ctx context.Context, pageNum, pageSize int) ([]*Review, error) {
	// 添加超时控制
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var reviews []*Review
	err := DB.WithContext(ctxWithTimeout).Table(constants.ReviewTableName).
		Where("target_type = ? AND status = ?", "resource", "pending").
		Order("created_at desc").
		Offset((pageNum - 1) * pageSize).
		Limit(pageSize).
		Find(&reviews).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询资源举报列表失败: "+err.Error())
	}
	return reviews, nil
}

func CountPendingResourceReviews(ctx context.Context) (int64, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var total int64
	err := DB.WithContext(ctxWithTimeout).Table(constants.ReviewTableName).
		Where("target_type = ? AND status = ?", "resource", "pending").
		Count(&total).Error
	if err != nil {
		return 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, "统计资源举报数量失败: "+err.Error())
	}
	return total, nil
}

// AuditResourceReview 审核资源举报记录
func AuditResourceReview(ctx context.Context, reviewID, reviewerID int64, action string) error {
	// 开始事务
	tx := DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 查询举报记录
	var review Review
	if err := tx.Table(constants.ReviewTableName).Where("review_id = ?", reviewID).First(&review).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return errno.NewErrNo(errno.InternalDatabaseErrorCode, "记录未找到")
		}
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询举报记录失败: "+err.Error())
	}

	// 校验举报目标类型
	if review.TargetType != "resource" {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "举报类型不匹配")
	}

	// 计算更新后的状态
	var newStatus string
	switch action {
	case "approve":
		newStatus = "approved"
	case "reject":
		newStatus = "rejected"
	default:
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "操作类型无效")
	}

	// 更新举报状态
	now := time.Now()
	if err := tx.Table(constants.ReviewTableName).Where("review_id = ?", reviewID).Updates(map[string]interface{}{
		"status":      newStatus,
		"reviewer_id": reviewerID,
		"reviewed_at": now,
	}).Error; err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新举报状态失败: "+err.Error())
	}

	// 根据审核结果更新资源状态
	if newStatus == "approved" {
		if err := tx.Table(constants.ResourceTableName).Where("resource_id = ?", review.TargetID).
			Update("status", "banned").Error; err != nil {
			tx.Rollback()
			return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新资源状态失败: "+err.Error())
		}
	} else if newStatus == "rejected" {
		if err := tx.Table(constants.ResourceTableName).Where("resource_id = ?", review.TargetID).
			Update("status", "normal").Error; err != nil {
			tx.Rollback()
			return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新资源状态失败: "+err.Error())
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "提交审核事务失败: "+err.Error())
	}
	return nil
}

func AuditResourceCommentReview(ctx context.Context, reviewID, reviewerID int64, action string) error {
	tx := DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var review Review
	if err := tx.Table(constants.ReviewTableName).Where("review_id = ?", reviewID).First(&review).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return errno.NewErrNo(errno.InternalDatabaseErrorCode, "记录未找到")
		}
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询举报记录失败: "+err.Error())
	}

	if review.TargetType != "comment" {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "举报类型不匹配")
	}

	var newStatus string
	switch action {
	case "approve":
		newStatus = "approved"
	case "reject":
		newStatus = "rejected"
	default:
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "操作类型无效")
	}

	now := time.Now()
	if err := tx.Table(constants.ReviewTableName).Where("review_id = ?", reviewID).Updates(map[string]interface{}{
		"status":      newStatus,
		"reviewer_id": reviewerID,
		"reviewed_at": now,
	}).Error; err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新举报状态失败: "+err.Error())
	}

	if newStatus == "approved" {
		if err := tx.Table(constants.ResourceCommentTableName).Where("comment_id = ?", review.TargetID).
			Updates(map[string]interface{}{
				"status":     "deleted_by_admin",
				"is_visible": false,
			}).Error; err != nil {
			tx.Rollback()
			return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新资源评论状态失败: "+err.Error())
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "提交审核事务失败: "+err.Error())
	}
	return nil
}

func GetPendingCommentReviews(ctx context.Context, pageNum, pageSize int) ([]*Review, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var reviews []*Review
	err := DB.WithContext(ctxWithTimeout).Table(constants.ReviewTableName).
		Where("target_type = ? AND status = ?", "comment", "pending").
		Order("created_at desc").
		Offset((pageNum - 1) * pageSize).
		Limit(pageSize).
		Find(&reviews).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询待审核资源评论的审核记录失败: "+err.Error())
	}
	return reviews, nil
}
