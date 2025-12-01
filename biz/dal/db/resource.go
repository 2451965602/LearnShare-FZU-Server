package db

import (
	"LearnShare/pkg/constants"
	"LearnShare/pkg/errno"
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func SearchResources(ctx context.Context, keyword *string, tagID, courseID *int64, sortBy *string, pageNum, pageSize int) ([]*Resource, int64, error) {
	// 添加超时控制
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var resources []*Resource
	var total int64

	db := DB.WithContext(ctxWithTimeout).Table(constants.ResourceTableName)

	if keyword != nil && *keyword != "" {
		db = db.Where("resource_name LIKE ? OR description LIKE ?", "%"+*keyword+"%", "%"+*keyword+"%")
	}

	if courseID != nil {
		db = db.Where("course_id = ?", *courseID)
	}

	if tagID != nil {
		db = db.Joins("JOIN "+constants.ResourceTagMappingTableName+" ON "+constants.ResourceTagMappingTableName+".resource_id = "+constants.ResourceTableName+".resource_id").
			Where(constants.ResourceTagMappingTableName+".tag_id = ?", *tagID)
	}

	switch {
	case sortBy != nil && *sortBy == "hot":
		db = db.Order("download_count desc")
	case sortBy != nil && *sortBy == "rating":
		db = db.Order("average_rating desc")
	default:
		db = db.Order("created_at desc")
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, "统计资源数量失败: "+err.Error())
	}

	err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Preload("Tags").Find(&resources).Error
	if err != nil {
		return nil, 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询资源列表失败: "+err.Error())
	}

	return resources, total, nil
}

// GetResourceByID 根据资源ID获取单个资源信息
func GetResourceByID(ctx context.Context, resourceID int64) (*Resource, error) {
	var resource Resource

	err := DB.WithContext(ctx).Table(constants.ResourceTableName).
		Preload("Tags").
		Where("resource_id = ?", resourceID).
		First(&resource).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "记录未找到")
		}
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "根据ID查询资源失败: "+err.Error())
	}

	return &resource, nil
}

// GetResourceComments 获取资源评论列表
func GetResourceComments(ctx context.Context, resourceID int64, sortBy *string, pageNum, pageSize int) ([]*ResourceCommentWithUser, int64, error) {
	// 添加超时控制
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 先统计总数
	var total int64
	base := DB.WithContext(ctxWithTimeout).Table(constants.ResourceCommentTableName).Model(&ResourceComment{}).
		Where("resource_id = ?", resourceID).
		Where("is_visible = ?", true).
		Where("status = ?", "normal")

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, "统计资源评论数量失败: "+err.Error())
	}

	// 使用联表查询，一次性获取评论与用户信息

	var rows []ResourceCommentrow

	query := DB.WithContext(ctxWithTimeout).Table(constants.ResourceCommentTableName+" AS c").Select(
		"c.comment_id, c.user_id, c.resource_id, c.content, c.parent_id, c.likes, c.is_visible, c.status, c.created_at, "+
			"u.user_id AS u_user_id, u.username AS u_username, u.email AS u_email, u.college_id AS u_college_id, u.major_id AS u_major_id, u.avatar_url AS u_avatar_url, u.reputation_score AS u_reputation_score, u.role_id AS u_role_id, u.status AS u_status, u.created_at AS u_created_at, u.updated_at AS u_updated_at",
	).Joins("LEFT JOIN "+constants.UserTableName+" u ON c.user_id = u.user_id").Where("c.resource_id = ? AND c.is_visible = ? AND c.status = ?", resourceID, true, "normal")

	// 排序
	switch {
	case sortBy != nil && *sortBy == "hottest":
		query = query.Order("c.likes DESC, c.created_at DESC")
	default:
		query = query.Order("c.created_at DESC")
	}

	err := query.Offset((pageNum - 1) * pageSize).Limit(pageSize).Scan(&rows).Error
	if err != nil {
		return nil, 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询资源评论列表失败: "+err.Error())
	}

	if len(rows) == 0 {
		return []*ResourceCommentWithUser{}, total, nil
	}

	// 映射为目标结构
	result := make([]*ResourceCommentWithUser, 0, len(rows))
	for _, r := range rows {
		var user User
		if r.UUserID == nil {
			user = User{UserID: 0, Username: "未知用户"}
		} else {
			user = User{
				UserID:          *r.UUserID,
				Username:        derefString(r.UUsername),
				PasswordHash:    derefString(r.UPasswordHash),
				Email:           derefString(r.UEmail),
				CollegeID:       r.UCollegeID,
				MajorID:         r.UMajorID,
				AvatarURL:       r.UAvatarURL,
				ReputationScore: derefInt64(r.UReputationScore),
				RoleID:          derefInt64(r.URoleID),
				Status:          derefString(r.UStatus),
			}
			if r.UCreatedAt != nil {
				user.CreatedAt = *r.UCreatedAt
			}
			if r.UUpdatedAt != nil {
				user.UpdatedAt = *r.UUpdatedAt
			}
		}

		cc := &ResourceCommentWithUser{
			CommentID:  r.CommentID,
			UserID:     r.UserID,
			ResourceID: r.ResourceID,
			Content:    r.Content,
			ParentID:   r.ParentID,
			Likes:      r.Likes,
			IsVisible:  r.IsVisible,
			Status:     r.Status,
			CreatedAt:  r.CreatedAt,
			User:       user,
		}
		result = append(result, cc)
	}

	return result, total, nil
}

// SubmitResourceRating 提交资源评分
func SubmitResourceRating(ctx context.Context, userID, resourceID int64, recommendation float64) (*ResourceRating, error) {
	// 开始事务
	tx := DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查是否已经评分过
	var existingRating ResourceRating
	err := tx.Table(constants.ResourceRatingTableName).Where("user_id = ? AND resource_id = ?", userID, resourceID).Find(&existingRating).Error
	if err != nil {
		tx.Rollback()
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询资源评分记录失败: "+err.Error())
	}

	var rating *ResourceRating

	if existingRating.RatingID > 0 {
		// 更新现有评分
		existingRating.Recommendation = recommendation
		existingRating.IsVisible = true // 确保在重新评分时，记录是可见的
		err = tx.Table(constants.ResourceRatingTableName).Save(&existingRating).Error
		rating = &existingRating
	} else {
		// 创建新评分
		rating = &ResourceRating{
			UserID:         userID,
			ResourceID:     resourceID,
			Recommendation: recommendation,
			IsVisible:      true,
		}
		err = tx.Table(constants.ResourceRatingTableName).Create(rating).Error
	}

	if err != nil {
		tx.Rollback()
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "提交资源评分失败: "+err.Error())
	}

	// 重新计算资源的平均评分
	var avgResult struct {
		AverageRating float64 `gorm:"column:average_rating"`
		RatingCount   int64   `gorm:"column:rating_count"`
	}

	err = tx.Table(constants.ResourceRatingTableName).
		Select("AVG(recommendation) as average_rating, COUNT(*) as rating_count").
		Where("resource_id = ? AND is_visible = ?", resourceID, true).
		Scan(&avgResult).Error

	if err != nil {
		tx.Rollback()
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "计算资源平均评分失败: "+err.Error())
	}

	// 更新资源的评分信息
	err = tx.Table(constants.ResourceTableName).
		Where("resource_id = ?", resourceID).
		Updates(map[string]interface{}{
			"average_rating": avgResult.AverageRating,
			"rating_count":   avgResult.RatingCount,
		}).Error

	if err != nil {
		tx.Rollback()
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新资源评分信息失败: "+err.Error())
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "提交评分事务失败: "+err.Error())
	}

	return rating, nil
}

// SubmitResourceComment 提交资源评论
func SubmitResourceComment(ctx context.Context, userID, resourceID int64, content string, parentID *int64) (*ResourceComment, error) {
	// 开始事务
	tx := DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建评论
	comment := &ResourceComment{
		UserID:     userID,
		ResourceID: resourceID,
		Content:    content,
		ParentID:   parentID,
		Likes:      0,
		IsVisible:  true,
		Status:     "normal",
	}

	// 保存评论
	err := tx.Table(constants.ResourceCommentTableName).Create(comment).Error
	if err != nil {
		tx.Rollback()
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "保存资源评论失败: "+err.Error())
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "提交评论事务失败: "+err.Error())
	}

	// 预加载用户信息
	err2 := DB.WithContext(ctx).Table(constants.ResourceCommentTableName).Preload("User").First(comment, comment.CommentID).Error
	if err2 != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "预加载评论用户信息失败: "+err2.Error())
	}

	return comment, nil
}

// SubmitResourceCommentAsync 异步提交资源评论
func SubmitResourceCommentAsync(ctx context.Context, userID, resourceID int64, content string, parentID *int64) chan struct {
	Comment *ResourceComment
	Err     error
} {
	resultChan := make(chan struct {
		Comment *ResourceComment
		Err     error
	}, 1)

	go func() {
		comment, err := SubmitResourceComment(ctx, userID, resourceID, content, parentID)
		resultChan <- struct {
			Comment *ResourceComment
			Err     error
		}{Comment: comment, Err: err}
		close(resultChan)
	}()

	return resultChan
}

// DeleteResourceRating 删除资源评分
func DeleteResourceRating(ctx context.Context, ratingID, userID int64) error {
	// 开始事务
	tx := DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 查询评分记录，确保用户只能删除自己的评分
	var rating ResourceRating
	err := tx.Table(constants.ResourceRatingTableName).Where("rating_id = ? AND user_id = ?", ratingID, userID).First(&rating).Error
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.NewErrNo(errno.InternalDatabaseErrorCode, "未找到评分记录或无权删除")
		}
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询评分记录失败: "+err.Error())
	}

	// 直接从数据库中删除评分记录
	err = tx.Table(constants.ResourceRatingTableName).Delete(&rating).Error
	if err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "删除评分失败: "+err.Error())
	}

	// 重新计算资源的平均评分
	var avgResult struct {
		AverageRating float64 `gorm:"column:average_rating"`
		RatingCount   int64   `gorm:"column:rating_count"`
	}

	err = tx.Table(constants.ResourceRatingTableName).
		Select("AVG(recommendation) as average_rating, COUNT(*) as rating_count").
		Where("resource_id = ?", rating.ResourceID).
		Scan(&avgResult).Error

	if err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "计算资源平均评分失败: "+err.Error())
	}

	// 更新资源的评分信息
	err = tx.Table(constants.ResourceTableName).
		Where("resource_id = ?", rating.ResourceID).
		Updates(map[string]interface{}{
			"average_rating": avgResult.AverageRating,
			"rating_count":   avgResult.RatingCount,
		}).Error

	if err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新资源评分信息失败: "+err.Error())
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "提交删除评分事务失败: "+err.Error())
	}

	return nil
}

// DeleteResourceRatingAsync 异步删除资源评分
func DeleteResourceRatingAsync(ctx context.Context, ratingID, userID int64) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return DeleteResourceRating(ctx, ratingID, userID)
	})
}

// DeleteResourceComment 删除资源评论
func DeleteResourceComment(ctx context.Context, commentID, userID int64) error {
	// 开始事务
	tx := DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 直接删除评论，确保用户只能删除自己的评论
	result := tx.Table(constants.ResourceCommentTableName).Where("comment_id = ? AND user_id = ?", commentID, userID).Delete(&ResourceComment{})
	if result.Error != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "删除评论失败: "+result.Error.Error())
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "未找到评论或无权删除")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "提交删除评论事务失败: "+err.Error())
	}

	return nil
}

// DeleteResourceCommentAsync 异步删除资源评论
func DeleteResourceCommentAsync(ctx context.Context, commentID, userID int64) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return DeleteResourceComment(ctx, commentID, userID)
	})
}

// AdminDeleteResourceComment 管理员删除资源评论（不限制用户）
func AdminDeleteResourceComment(ctx context.Context, commentID int64) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx := DB.WithContext(ctxWithTimeout).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var comment ResourceComment
	if err := tx.Table(constants.ResourceCommentTableName).Where("comment_id = ?", commentID).First(&comment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return errno.NewErrNo(errno.InternalDatabaseErrorCode, "记录未找到")
		}
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询评论失败: "+err.Error())
	}

	if err := tx.Table(constants.ResourceCommentTableName).Where("comment_id = ?", commentID).Delete(&ResourceComment{}).Error; err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "删除评论失败: "+err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "提交删除事务失败: "+err.Error())
	}
	return nil
}

// AdminDeleteResourceRating 管理员删除资源评分（不限制用户）并重算平均分
func AdminDeleteResourceRating(ctx context.Context, ratingID int64) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx := DB.WithContext(ctxWithTimeout).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var rating ResourceRating
	if err := tx.Table(constants.ResourceRatingTableName).Where("rating_id = ?", ratingID).First(&rating).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return errno.NewErrNo(errno.InternalDatabaseErrorCode, "记录未找到")
		}
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询评分失败: "+err.Error())
	}

	if err := tx.Table(constants.ResourceRatingTableName).Where("rating_id = ?", ratingID).Delete(&ResourceRating{}).Error; err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "删除评分失败: "+err.Error())
	}

	var avgResult struct {
		AverageRating float64 `gorm:"column:average_rating"`
		RatingCount   int64   `gorm:"column:rating_count"`
	}
	if err := tx.Table(constants.ResourceRatingTableName).
		Select("AVG(recommendation) as average_rating, COUNT(*) as rating_count").
		Where("resource_id = ? AND is_visible = ?", rating.ResourceID, true).
		Scan(&avgResult).Error; err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "计算资源平均评分失败: "+err.Error())
	}

	if err := tx.Table(constants.ResourceTableName).
		Where("resource_id = ?", rating.ResourceID).
		Updates(map[string]interface{}{
			"average_rating": avgResult.AverageRating,
			"rating_count":   avgResult.RatingCount,
		}).Error; err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新资源评分信息失败: "+err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "提交删除评分事务失败: "+err.Error())
	}
	return nil
}

// AdminDeleteResource 管理员硬删除资源，并清理关联引用
func AdminDeleteResource(ctx context.Context, resourceID int64) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx := DB.WithContext(ctxWithTimeout).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var res Resource
	if err := tx.Table(constants.ResourceTableName).Where("resource_id = ?", resourceID).First(&res).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return errno.NewErrNo(errno.InternalDatabaseErrorCode, "记录未找到")
		}
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询资源失败: "+err.Error())
	}

	if err := tx.Table(constants.ResourceTableName).Where("resource_id = ?", resourceID).Delete(&Resource{}).Error; err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "删除资源失败: "+err.Error())
	}

	// 清理 favorites 的孤儿记录（非外键约束，需要手动清理）
	if err := tx.Table(constants.FavoriteTableName).Where("target_type = ? AND target_id = ?", "resource", resourceID).Delete(nil).Error; err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, fmt.Sprintf("清理收藏引用失败: %v", err))
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "提交删除事务失败: "+err.Error())
	}

	return nil
}

// CreateReview 创建一个新的举报（审核）
func CreateReview(ctx context.Context, creatorID int64, targetID int64, targetType, reason string) error {
	review := &Review{
		TargetID:   targetID,
		TargetType: targetType,
		Reason:     reason,
		ReporterID: creatorID, // 举报者
		ReviewerID: nil,       // 审核人暂为空，待审核分配
	}

	result := DB.WithContext(ctx).Table(constants.ReviewTableName).Create(review)
	if result.Error != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "创建举报失败: "+result.Error.Error())
	}

	return nil
}

// CreateReviewAsync 异步创建举报
func CreateReviewAsync(ctx context.Context, creatorID int64, targetID int64, targetType, reason string) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return CreateReview(ctx, creatorID, targetID, targetType, reason)
	})
}

// GetOrCreateTag 获取或创建单个标签
func GetOrCreateTag(ctx context.Context, tagName string) (*ResourceTag, error) {
	var tag ResourceTag
	err := DB.WithContext(ctx).Table("tags").Where("tag_name = ?", tagName).First(&tag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tag = ResourceTag{TagName: tagName}
			if e := DB.WithContext(ctx).Table("tags").Create(&tag).Error; e != nil {
				return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "创建标签失败: "+e.Error())
			}
		} else {
			return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询标签失败: "+err.Error())
		}
	}
	return &tag, nil
}

// GetOrCreateTagsBatch 批量获取或创建标签（优化：减少数据库往返）
func GetOrCreateTagsBatch(ctx context.Context, tagNames []string) ([]*ResourceTag, error) {
	if len(tagNames) == 0 {
		return []*ResourceTag{}, nil
	}

	// 先去重
	uniqueNames := make(map[string]bool)
	for _, name := range tagNames {
		if trimmed := name; trimmed != "" {
			uniqueNames[trimmed] = true
		}
	}

	var uniqueList []string
	for name := range uniqueNames {
		uniqueList = append(uniqueList, name)
	}

	// 查询已存在的标签
	var existingTags []ResourceTag
	err := DB.WithContext(ctx).Table("tags").Where("tag_name IN (?)", uniqueList).Find(&existingTags).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "批量查询标签失败: "+err.Error())
	}

	// 建立已存在标签的映射
	existingMap := make(map[string]*ResourceTag)
	for i := range existingTags {
		existingMap[existingTags[i].TagName] = &existingTags[i]
	}

	// 找出需要创建的新标签
	var newTags []*ResourceTag
	for _, name := range uniqueList {
		if _, exists := existingMap[name]; !exists {
			newTags = append(newTags, &ResourceTag{TagName: name})
		}
	}

	// 批量创建新标签
	if len(newTags) > 0 {
		if err := DB.WithContext(ctx).Table("tags").CreateInBatches(newTags, 50).Error; err != nil {
			// 如果是重复键错误，说明有其他并发操作，可以忽略
			if !errors.Is(err, gorm.ErrDuplicatedKey) {
				return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "批量创建标签失败: "+err.Error())
			}
		}
	}

	// 重新查询所有标签以确保数据完整性
	var allTags []ResourceTag
	err = DB.WithContext(ctx).Table("tags").Where("tag_name IN (?)", uniqueList).Find(&allTags).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "重新查询标签失败: "+err.Error())
	}

	// 转换为指针数组
	result := make([]*ResourceTag, len(allTags))
	for i := range allTags {
		result[i] = &allTags[i]
	}

	return result, nil
}

// LinkResourceTag 关联单个资源标签
func LinkResourceTag(ctx context.Context, resourceID, tagID int64) error {
	mapping := ResourceTagMapping{ResourceID: resourceID, TagID: tagID}
	if err := DB.WithContext(ctx).Table(constants.ResourceTagMappingTableName).Create(&mapping).Error; err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "关联资源标签失败: "+err.Error())
	}
	return nil
}

// LinkResourceTagsBatch 批量关联资源标签（优化：减少数据库往返）
func LinkResourceTagsBatch(ctx context.Context, resourceID int64, tagIDs []int64) error {
	if len(tagIDs) == 0 {
		return nil
	}

	mappings := make([]ResourceTagMapping, 0, len(tagIDs))
	for _, tagID := range tagIDs {
		if tagID > 0 {
			mappings = append(mappings, ResourceTagMapping{
				ResourceID: resourceID,
				TagID:      tagID,
			})
		}
	}

	if len(mappings) == 0 {
		return nil
	}

	// 批量插入，忽略重复键错误
	err := DB.WithContext(ctx).Table(constants.ResourceTagMappingTableName).CreateInBatches(mappings, 100).Error
	if err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "批量关联资源标签失败: "+err.Error())
	}

	return nil
}

func ReactResourceComment(ctx context.Context, userID, commentID int64, action string) error {
	tx := DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var comment ResourceComment
	if err := tx.Table(constants.ResourceCommentTableName).Where("comment_id = ? AND is_visible = ? AND status = ?", commentID, true, "normal").First(&comment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return errno.ResourceInvalidCommentError
		}
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询评论失败: "+err.Error())
	}

	var existing ResourceCommentReaction
	if err := tx.Table(constants.ResourceCommentReactionTableName).Where("user_id = ? AND comment_id = ?", userID, commentID).First(&existing).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询评论反应失败: "+err.Error())
	}

	switch action {
	case "like":
		if existing.ReactionID == 0 {
			rec := &ResourceCommentReaction{UserID: userID, CommentID: commentID, Reaction: "like"}
			if err := tx.Table(constants.ResourceCommentReactionTableName).Create(rec).Error; err != nil {
				tx.Rollback()
				return errno.NewErrNo(errno.InternalDatabaseErrorCode, "创建评论点赞记录失败: "+err.Error())
			}
			if err := tx.Table(constants.ResourceCommentTableName).Where("comment_id = ?", commentID).Update("likes", gorm.Expr("likes + 1")).Error; err != nil {
				tx.Rollback()
				return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新点赞数失败: "+err.Error())
			}
		} else if existing.Reaction == "dislike" {
			if err := tx.Table(constants.ResourceCommentReactionTableName).Where("reaction_id = ?", existing.ReactionID).Update("reaction", "like").Error; err != nil {
				tx.Rollback()
				return errno.NewErrNo(errno.InternalDatabaseErrorCode, "切换为点赞失败: "+err.Error())
			}
			if err := tx.Table(constants.ResourceCommentTableName).Where("comment_id = ?", commentID).Update("likes", gorm.Expr("likes + 1")).Error; err != nil {
				tx.Rollback()
				return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新点赞数失败: "+err.Error())
			}
		}
	case "dislike":
		if existing.ReactionID == 0 {
			rec := &ResourceCommentReaction{UserID: userID, CommentID: commentID, Reaction: "dislike"}
			if err := tx.Table(constants.ResourceCommentReactionTableName).Create(rec).Error; err != nil {
				tx.Rollback()
				return errno.NewErrNo(errno.InternalDatabaseErrorCode, "创建评论点踩记录失败: "+err.Error())
			}
		} else if existing.Reaction == "like" {
			if err := tx.Table(constants.ResourceCommentReactionTableName).Where("reaction_id = ?", existing.ReactionID).Update("reaction", "dislike").Error; err != nil {
				tx.Rollback()
				return errno.NewErrNo(errno.InternalDatabaseErrorCode, "切换为点踩失败: "+err.Error())
			}
			if err := tx.Table(constants.ResourceCommentTableName).Where("comment_id = ? AND likes > 0", commentID).Update("likes", gorm.Expr("likes - 1")).Error; err != nil {
				tx.Rollback()
				return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新点赞数失败: "+err.Error())
			}
		}
	case "cancel_like":
		if existing.ReactionID != 0 && existing.Reaction == "like" {
			if err := tx.Table(constants.ResourceCommentReactionTableName).Where("reaction_id = ?", existing.ReactionID).Delete(&ResourceCommentReaction{}).Error; err != nil {
				tx.Rollback()
				return errno.NewErrNo(errno.InternalDatabaseErrorCode, "撤销点赞失败: "+err.Error())
			}
			if err := tx.Table(constants.ResourceCommentTableName).Where("comment_id = ? AND likes > 0", commentID).Update("likes", gorm.Expr("likes - 1")).Error; err != nil {
				tx.Rollback()
				return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新点赞数失败: "+err.Error())
			}
		}
	case "cancel_dislike":
		if existing.ReactionID != 0 && existing.Reaction == "dislike" {
			if err := tx.Table(constants.ResourceCommentReactionTableName).Where("reaction_id = ?", existing.ReactionID).Delete(&ResourceCommentReaction{}).Error; err != nil {
				tx.Rollback()
				return errno.NewErrNo(errno.InternalDatabaseErrorCode, "撤销点踩失败: "+err.Error())
			}
		}
	default:
		tx.Rollback()
		return errno.ParamVerifyError.WithMessage("操作类型无效")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "提交评论反应事务失败: "+err.Error())
	}
	return nil
}

func ReactResourceCommentAsync(ctx context.Context, userID, commentID int64, action string) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return ReactResourceComment(ctx, userID, commentID, action)
	})
}
