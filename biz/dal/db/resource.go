package db

import (
	"context"
)

func SearchResources(ctx context.Context, keyword *string, tagID, courseID *int64, sortBy *string, pageNum, pageSize int) ([]*Resource, int64, error) {
	var resources []*Resource
	var total int64

	// 参数验证
	if pageNum <= 0 || pageSize <= 0 {
		return []*Resource{}, 0, nil
	}

	db := DB.WithContext(ctx)

	if keyword != nil && *keyword != "" {
		db = db.Where("title LIKE ? OR description LIKE ?", "%"+*keyword+"%", "%"+*keyword+"%")
	}

	if courseID != nil {
		db = db.Where("course_id = ?", *courseID)
	}

	if tagID != nil {
		db = db.Joins("JOIN resource_tag_mapping ON resource_tag_mapping.resource_id = resource.resource_id").
			Where("resource_tag_mapping.tag_id = ?", *tagID)
	}

	switch {
	case sortBy != nil && *sortBy == "hot":
		db = db.Order("download_count desc")
	case sortBy != nil && *sortBy == "rating":
		db = db.Order("average_rating desc")
	default:
		db = db.Order("created_at desc")
	}

	err := db.Model(&Resource{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Preload("Tags").Find(&resources).Error
	if err != nil {
		return nil, 0, err
	}

	return resources, total, nil
}

// GetResourceByID 根据资源ID获取单个资源信息
func GetResourceByID(ctx context.Context, resourceID int64) (*Resource, error) {
	var resource Resource

	err := DB.WithContext(ctx).
		Preload("Tags").
		Where("resource_id = ?", resourceID).
		First(&resource).Error

	if err != nil {
		return nil, err
	}

	return &resource, nil
}

// GetResourceComments 获取资源评论列表
func GetResourceComments(ctx context.Context, resourceID int64, sortBy *string, pageNum, pageSize int) ([]*ResourceComment, int64, error) {
	var comments []*ResourceComment
	var total int64

	// 参数验证
	if pageNum <= 0 || pageSize <= 0 {
		return []*ResourceComment{}, 0, nil
	}

	db := DB.WithContext(ctx).
		Preload("User").
		Where("resource_id = ?", resourceID).
		Where("is_visible = ?", true).
		Where("status = ?", "normal")

	// 根据排序参数进行排序
	if sortBy != nil {
		switch *sortBy {
		case "latest":
			db = db.Order("created_at DESC")
		case "hottest":
			db = db.Order("likes DESC, created_at DESC")
		default:
			db = db.Order("created_at DESC")
		}
	} else {
		db = db.Order("created_at DESC")
	}

	// 获取总数
	err := db.Model(&ResourceComment{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = db.Offset((pageNum - 1) * pageSize).
		Limit(pageSize).
		Find(&comments).Error
	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
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
	err := tx.Where("user_id = ? AND resource_id = ?", userID, resourceID).First(&existingRating).Error
	
	var rating *ResourceRating
	
	if err == nil {
		// 更新现有评分
		existingRating.Recommendation = recommendation
		err = tx.Save(&existingRating).Error
		rating = &existingRating
	} else {
		// 创建新评分
		rating = &ResourceRating{
			UserID:         userID,
			ResourceID:     resourceID,
			Recommendation: recommendation,
			IsVisible:      true,
		}
		err = tx.Create(rating).Error
	}
	
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 重新计算资源的平均评分
	var avgResult struct {
		AverageRating float64 `gorm:"column:average_rating"`
		RatingCount   int64   `gorm:"column:rating_count"`
	}
	
	err = tx.Model(&ResourceRating{}).
		Select("AVG(recommendation) as average_rating, COUNT(*) as rating_count").
		Where("resource_id = ? AND is_visible = ?", resourceID, true).
		Scan(&avgResult).Error
	
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 更新资源的评分信息
	err = tx.Model(&Resource{}).
		Where("resource_id = ?", resourceID).
		Updates(map[string]interface{}{
			"average_rating": avgResult.AverageRating,
			"rating_count":   avgResult.RatingCount,
		}).Error
	
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
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
	err := tx.Create(comment).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 预加载用户信息
	err2 := DB.WithContext(ctx).Preload("User").First(comment, comment.CommentID).Error
	if err2 != nil {
		return nil, err2
	}

	return comment, nil
}
