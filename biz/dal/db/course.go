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

func CreateCourse(ctx context.Context, courseName string, teacherID, majorID int64, credit float64, grade, description string) error {
	course := &Course{
		CourseName:  courseName,
		TeacherID:   teacherID,
		Credit:      credit,
		MajorID:     majorID,
		Grade:       grade,
		Description: &description,
	}

	err := DB.WithContext(ctx).Table(constants.CourseTableName).Create(course).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "创建课程失败: "+err.Error())
	}
	return nil
}

// CreateCourseAsync 异步创建课程
func CreateCourseAsync(ctx context.Context, courseName string, teacherID, majorID int64, credit float64, grade, description string) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return CreateCourse(ctx, courseName, teacherID, majorID, credit, grade, description)
	})
}

func UpdateCourse(ctx context.Context, courseID int64, updates map[string]interface{}) error {
	err := DB.WithContext(ctx).Table(constants.CourseTableName).Where("course_id = ?", courseID).Updates(updates).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新课程失败: "+err.Error())
	}
	return nil
}

// UpdateCourseAsync 异步更新课程
func UpdateCourseAsync(ctx context.Context, courseID int64, updates map[string]interface{}) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return UpdateCourse(ctx, courseID, updates)
	})
}

func DeleteCourse(ctx context.Context, courseID int64) error {
	err := DB.WithContext(ctx).Table(constants.CourseTableName).Where("course_id = ?", courseID).Delete(&Course{}).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "删除课程失败: "+err.Error())
	}
	return nil
}

// DeleteCourseAsync 异步删除课程
func DeleteCourseAsync(ctx context.Context, courseID int64) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return DeleteCourse(ctx, courseID)
	})
}

func GetCourseByID(ctx context.Context, courseID int64) (*Course, error) {
	var course Course
	err := DB.WithContext(ctx).Table(constants.CourseTableName).Where("course_id = ?", courseID).First(&course).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询课程失败: "+err.Error())
	}
	return &course, nil
}

func GetCoursesByTeacherID(ctx context.Context, teacherID int64, pageSize, pageNum int) ([]*Course, error) {
	var courses []*Course
	err := DB.WithContext(ctx).Table(constants.CourseTableName).Where("teacher_id = ?", teacherID).Limit(pageSize).Offset(pageSize * (pageNum - 1)).Find(&courses).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询教师课程失败: "+err.Error())
	}
	return courses, nil
}

func GetCoursesByMajorID(ctx context.Context, majorID int64) ([]*Course, error) {
	var courses []*Course
	err := DB.WithContext(ctx).Table(constants.CourseTableName).Where("major_id = ?", majorID).Find(&courses).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询专业课程失败: "+err.Error())
	}
	return courses, nil
}

func SearchCourses(ctx context.Context, keywords string, grade string, pageNum, pageSize int) ([]*Course, error) {
	var courses []*Course // 声明courses变量

	query := DB.WithContext(ctx).Table(constants.CourseTableName)

	if keywords != "" {
		query = query.Where("course_name LIKE ?", "%"+keywords+"%")
	}
	if grade != "" {
		query = query.Where("grade = ?", grade)
	}

	err := query.Limit(pageSize).Offset(pageSize * (pageNum - 1)).Find(&courses).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "搜索课程失败: "+err.Error())
	}
	return courses, nil
}

func SubmitCourseRating(ctx context.Context, rating *CourseRating) (*CourseRating, error) {
	tx := DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var existing CourseRating
	err := tx.Table(constants.CourseRatingTableName).
		Where("user_id = ? AND course_id = ?", rating.UserID, rating.CourseID).
		First(&existing).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询课程评分记录失败: "+err.Error())
	}

	if err == nil {
		// 更新
		existing.Recommendation = rating.Recommendation
		existing.Difficulty = rating.Difficulty
		existing.Workload = rating.Workload
		existing.Usefulness = rating.Usefulness
		existing.IsVisible = true
		existing.UpdatedAt = time.Now()

		if err := tx.Table(constants.CourseRatingTableName).Where("rating_id = ?", rating.RatingID).Save(&existing).Error; err != nil {
			tx.Rollback()
			return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新课程评分失败: "+err.Error())
		}
		if err := tx.Commit().Error; err != nil {
			return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "提交评分事务失败: "+err.Error())
		}
		return &existing, nil
	}

	// 创建
	if err := tx.Table(constants.CourseRatingTableName).Create(rating).Error; err != nil {
		tx.Rollback()
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "提交课程评分失败: "+err.Error())
	}
	if err := tx.Commit().Error; err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "提交评分事务失败: "+err.Error())
	}
	return rating, nil
}

// SubmitCourseRatingAsync 异步提交课程评分
func SubmitCourseRatingAsync(ctx context.Context, rating *CourseRating) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		_, err := SubmitCourseRating(ctx, rating)
		return err
	})
}

func UpdateCourseRating(ctx context.Context, ratingID int64, updates map[string]interface{}) error {
	err := DB.WithContext(ctx).Table(constants.CourseRatingTableName).Where("rating_id = ?", ratingID).Updates(updates).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新课程评分失败: "+err.Error())
	}
	return nil
}

// UpdateCourseRatingAsync 异步更新课程评分
func UpdateCourseRatingAsync(ctx context.Context, ratingID int64, updates map[string]interface{}) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return UpdateCourseRating(ctx, ratingID, updates)
	})
}

func DeleteCourseRating(ctx context.Context, ratingID int64) error {
	err := DB.WithContext(ctx).Table(constants.CourseRatingTableName).Where("rating_id = ?", ratingID).Delete(&CourseRating{}).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "删除课程评分失败: "+err.Error())
	}
	return nil
}

// DeleteCourseRatingAsync 异步删除课程评分
func DeleteCourseRatingAsync(ctx context.Context, ratingID int64) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return DeleteCourseRating(ctx, ratingID)
	})
}

func GetCourseRatingByID(ctx context.Context, ratingID int64) (*CourseRating, error) {
	var rating CourseRating
	err := DB.WithContext(ctx).Table(constants.CourseRatingTableName).Where("rating_id = ?", ratingID).First(&rating).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询课程评分失败: "+err.Error())
	}
	return &rating, nil
}

func GetCourseRatingsByCourseID(ctx context.Context, courseID int64) ([]*CourseRating, error) {
	var ratings []*CourseRating
	err := DB.WithContext(ctx).Table(constants.CourseRatingTableName).Where("course_id = ? AND is_visible = ?", courseID, true).Order("created_at DESC").Find(&ratings).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询课程评分列表失败: "+err.Error())
	}
	return ratings, nil
}

func SubmitCourseComment(ctx context.Context, comment *CourseComment) error {
	err := DB.WithContext(ctx).Table(constants.CourseCommentTableName).Create(comment).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "提交课程评论失败: "+err.Error())
	}
	return nil
}

// SubmitCourseCommentAsync 异步提交课程评论
func SubmitCourseCommentAsync(ctx context.Context, comment *CourseComment) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return SubmitCourseComment(ctx, comment)
	})
}

func UpdateCourseComment(ctx context.Context, commentID int64, updates map[string]interface{}) error {
	err := DB.WithContext(ctx).Table(constants.CourseCommentTableName).Where("comment_id = ?", commentID).Updates(updates).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新课程评论失败: "+err.Error())
	}
	return nil
}

// UpdateCourseCommentAsync 异步更新课程评论
func UpdateCourseCommentAsync(ctx context.Context, commentID int64, updates map[string]interface{}) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return UpdateCourseComment(ctx, commentID, updates)
	})
}

func DeleteCourseComment(ctx context.Context, commentID int64) error {
	err := DB.WithContext(ctx).Table(constants.CourseCommentTableName).Where("comment_id = ?", commentID).Delete(&CourseComment{}).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "删除课程评论失败: "+err.Error())
	}
	return nil
}

// DeleteCourseCommentAsync 异步删除课程评论
func DeleteCourseCommentAsync(ctx context.Context, commentID int64) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return DeleteCourseComment(ctx, commentID)
	})
}

func GetCourseCommentByID(ctx context.Context, commentID int64) (*CourseComment, error) {
	var comment CourseComment
	err := DB.WithContext(ctx).Table(constants.CourseCommentTableName).Where("comment_id = ?", commentID).First(&comment).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询课程评论失败: "+err.Error())
	}
	return &comment, nil
}

func GetCourseCommentsByCourseID(ctx context.Context, courseID int64, sortBy string, pageNum, pageSize int) ([]*CourseCommentWithuser, error) {
	// 使用联表查询，一次性获取评论与用户信息

	var rows []CommentUserRow

	query := DB.WithContext(ctx).Table(constants.CourseCommentTableName+" AS c").Select(
		"c.comment_id, c.course_id, c.content, c.parent_id, c.is_visible, c.created_at, c.updated_at,c.likes,c.status,"+
			"u.user_id AS u_user_id, u.username AS u_username, u.email AS u_email, u.college_id AS u_college_id, u.major_id AS u_major_id, u.avatar_url AS u_avatar_url, u.reputation_score AS u_reputation_score, u.role_id AS u_role_id, u.status AS u_status, u.created_at AS u_created_at, u.updated_at AS u_updated_at,u.status AS u_status",
	).Joins("LEFT JOIN "+constants.UserTableName+" u ON c.user_id = u.user_id").Where("c.course_id = ? AND c.is_visible = ?", courseID, true)

	// 排序方式
	switch sortBy {
	case "latest":
		query = query.Order("c.created_at DESC")
	case "oldest":
		query = query.Order("c.created_at ASC")
	case "popular":
		query = query.Order("c.created_at DESC") // 简化处理
	default:
		query = query.Order("c.created_at DESC")
	}

	err := query.Limit(pageSize).Offset(pageSize * (pageNum - 1)).Scan(&rows).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询课程评论列表失败: "+err.Error())
	}

	if len(rows) == 0 {
		return []*CourseCommentWithuser{}, nil
	}

	// 将行映射为最终结构
	result := make([]*CourseCommentWithuser, 0, len(rows))
	for _, r := range rows {
		var user User
		if r.UserID == nil {
			user = User{UserID: 0, Username: "未知用户"}
		} else {
			user = User{
				UserID:          *r.UserID,
				Username:        derefString(r.Username),
				Email:           derefString(r.Email),
				CollegeID:       r.CollegeID,
				MajorID:         r.MajorID,
				AvatarURL:       r.AvatarURL,
				ReputationScore: derefInt64(r.ReputationScore),
				RoleID:          derefInt64(r.RoleID),
				Status:          r.Status,
			}
		}

		cc := &CourseCommentWithuser{
			CommentID: r.CommentID,
			CourseID:  r.CourseID,
			User:      user,
			Content:   r.Content,
			ParentID:  r.ParentID,
			IsVisible: r.IsVisible,
			Likes:     r.Likes,
			Status:    r.Status,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		}
		result = append(result, cc)
	}

	return result, nil
}

// 辅助函数，用于安全地解引用可空类型
func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefInt64(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}

// GetCourseResources 获取课程资源列表
func GetCourseResources(ctx context.Context, courseID int64, resourceType, status string, pageNum, pageSize int) ([]*Resource, error) {
	var resources []*Resource

	query := DB.WithContext(ctx).Table(constants.ResourceTableName).Where("course_id = ?", courseID)

	if resourceType != "" {
		query = query.Where("type = ?", resourceType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 使用 Preload 预加载 Tags，避免 N+1 查询问题
	err := query.Order("created_at DESC").Limit(pageSize).Offset(pageSize * (pageNum - 1)).Preload("Tags").Find(&resources).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询课程资源列表失败: "+err.Error())
	}

	return resources, nil
}

// CreateResource 创建资源
func CreateResource(ctx context.Context, resource *Resource) error {
	err := DB.WithContext(ctx).Table(constants.ResourceTableName).Create(resource).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "创建资源失败: "+err.Error())
	}
	return nil
}

// CreateResourceAsync 异步创建资源
func CreateResourceAsync(ctx context.Context, resource *Resource) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return CreateResource(ctx, resource)
	})
}

// UpdateResource 更新资源
func UpdateResource(ctx context.Context, resourceID int64, updates map[string]interface{}) error {
	err := DB.WithContext(ctx).Table(constants.ResourceTableName).Where("resource_id = ?", resourceID).Updates(updates).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新资源失败: "+err.Error())
	}
	return nil
}

// UpdateResourceAsync 异步更新资源
func UpdateResourceAsync(ctx context.Context, resourceID int64, updates map[string]interface{}) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return UpdateResource(ctx, resourceID, updates)
	})
}

// 删除资源
func DeleteResource(ctx context.Context, resourceID int64) error {
	err := DB.WithContext(ctx).Table(constants.ResourceTableName).Where("resource_id = ?", resourceID).Delete(&Resource{}).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "删除资源失败: "+err.Error())
	}
	return nil
}

// DeleteResourceAsync 异步删除资源
func DeleteResourceAsync(ctx context.Context, resourceID int64) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return DeleteResource(ctx, resourceID)
	})
}

// AdminDeleteCourse 管理员硬删除课程，并清理所有关联数据
func AdminDeleteCourse(ctx context.Context, courseID int64) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx := DB.WithContext(ctxWithTimeout).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 检查课程是否存在
	var course Course
	if err := tx.Table(constants.CourseTableName).Where("course_id = ?", courseID).First(&course).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return errno.NewErrNo(errno.InternalDatabaseErrorCode, "课程记录未找到")
		}
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询课程失败: "+err.Error())
	}

	// 2. 删除课程下的所有资源（假设资源表有 course_id 外键）
	if err := tx.Table(constants.ResourceTableName).Where("course_id = ?", courseID).Delete(&Resource{}).Error; err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "删除课程资源失败: "+err.Error())
	}

	// 3. 删除课程的所有评分
	if err := tx.Table(constants.CourseRatingTableName).Where("course_id = ?", courseID).Delete(&CourseRating{}).Error; err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "删除课程评分失败: "+err.Error())
	}

	// 4. 删除课程的所有评论
	if err := tx.Table(constants.CourseCommentTableName).Where("course_id = ?", courseID).Delete(&CourseComment{}).Error; err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "删除课程评论失败: "+err.Error())
	}

	// 5. 删除课程收藏记录（favorites 表，target_type = 'course'）
	if err := tx.Table(constants.FavoriteTableName).Where("target_type = ? AND target_id = ?", "course", courseID).Delete(nil).Error; err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, fmt.Sprintf("清理课程收藏引用失败: %v", err))
	}

	// 6. 最终删除课程本身
	if err := tx.Table(constants.CourseTableName).Where("course_id = ?", courseID).Delete(&Course{}).Error; err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "删除课程失败: "+err.Error())
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "提交删除课程事务失败: "+err.Error())
	}

	return nil
}

// AdminDeleteCourseComment 管理员删除课程评论（不限制用户）
func AdminDeleteCourseComment(ctx context.Context, commentID int64) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx := DB.WithContext(ctxWithTimeout).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var comment CourseComment
	if err := tx.Table(constants.CourseCommentTableName).Where("comment_id = ?", commentID).First(&comment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return errno.NewErrNo(errno.InternalDatabaseErrorCode, "课程评论记录未找到")
		}
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询课程评论失败: "+err.Error())
	}

	if err := tx.Table(constants.CourseCommentTableName).Where("comment_id = ?", commentID).Delete(&CourseComment{}).Error; err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "删除课程评论失败: "+err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "提交删除课程评论事务失败: "+err.Error())
	}
	return nil
}

// AdminDeleteCourseRating 管理员删除课程评分（不限制用户）并重算课程平均分
func AdminDeleteCourseRating(ctx context.Context, ratingID int64) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx := DB.WithContext(ctxWithTimeout).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var rating CourseRating
	if err := tx.Table(constants.CourseRatingTableName).Where("rating_id = ?", ratingID).First(&rating).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return errno.NewErrNo(errno.InternalDatabaseErrorCode, "课程评分记录未找到")
		}
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询课程评分失败: "+err.Error())
	}

	// 删除评分
	if err := tx.Table(constants.CourseRatingTableName).Where("rating_id = ?", ratingID).Delete(&CourseRating{}).Error; err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "删除课程评分失败: "+err.Error())
	}

	// 重新计算该课程的平均分和评分人数
	var avgResult struct {
		AverageRating float64 `gorm:"column:average_rating"`
		RatingCount   int64   `gorm:"column:rating_count"`
	}
	if err := tx.Table(constants.CourseRatingTableName).
		Select("AVG(rating) as average_rating, COUNT(*) as rating_count").
		Where("course_id = ? AND is_visible = ?", rating.CourseID, true).
		Scan(&avgResult).Error; err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "计算课程平均评分失败: "+err.Error())
	}

	// 更新课程表中的评分字段（假设 Course 表有 average_rating 和 rating_count 字段）
	if err := tx.Table(constants.CourseTableName).
		Where("course_id = ?", rating.CourseID).
		Updates(map[string]interface{}{
			"average_rating": avgResult.AverageRating,
			"rating_count":   avgResult.RatingCount,
		}).Error; err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新课程评分信息失败: "+err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "提交删除课程评分事务失败: "+err.Error())
	}
	return nil
}

// ReactCourseComment 处理课程评论的点赞/点踩/取消
func ReactCourseComment(ctx context.Context, userID, commentID int64, action string) error {
	tx := DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 验证评论是否存在且可见
	var comment CourseComment
	if err := tx.Table(constants.CourseCommentTableName).
		Where("comment_id = ? AND is_visible = ? AND status = ?", commentID, true, "normal").
		First(&comment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return errno.ResourceInvalidCommentError // 可复用或定义新的错误码
		}
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "查询课程评论失败: "+err.Error())
	}

	// 查询用户是否已有反应
	var existing CourseCommentReaction
	err := tx.Table(constants.CourseCommentReactionTableName).
		Where("user_id = ? AND comment_id = ?", userID, commentID).
		First(&existing).Error

	hasExisting := !errors.Is(err, gorm.ErrRecordNotFound)

	switch action {
	case "like":
		if !hasExisting {
			rec := &CourseCommentReaction{UserID: userID, CommentID: commentID, Reaction: "like"}
			if err := tx.Table(constants.CourseCommentReactionTableName).Create(rec).Error; err != nil {
				tx.Rollback()
				return errno.NewErrNo(errno.InternalDatabaseErrorCode, "创建课程评论点赞记录失败: "+err.Error())
			}
			if err := tx.Table(constants.CourseCommentTableName).Where("comment_id = ?", commentID).Update("likes", gorm.Expr("likes + 1")).Error; err != nil {
				tx.Rollback()
				return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新课程评论点赞数失败: "+err.Error())
			}
		} else if existing.Reaction == "dislike" {
			if err := tx.Table(constants.CourseCommentReactionTableName).Where("reaction_id = ?", existing.ReactionID).Update("reaction", "like").Error; err != nil {
				tx.Rollback()
				return errno.NewErrNo(errno.InternalDatabaseErrorCode, "切换为点赞失败: "+err.Error())
			}
			if err := tx.Table(constants.CourseCommentTableName).Where("comment_id = ?", commentID).Update("likes", gorm.Expr("likes + 1")).Error; err != nil {
				tx.Rollback()
				return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新点赞数失败: "+err.Error())
			}
		}
		// 若已是 like，则无操作

	case "dislike":
		if !hasExisting {
			rec := &CourseCommentReaction{UserID: userID, CommentID: commentID, Reaction: "dislike"}
			if err := tx.Table(constants.CourseCommentReactionTableName).Create(rec).Error; err != nil {
				tx.Rollback()
				return errno.NewErrNo(errno.InternalDatabaseErrorCode, "创建课程评论点踩记录失败: "+err.Error())
			}
		} else if existing.Reaction == "like" {
			if err := tx.Table(constants.CourseCommentReactionTableName).Where("reaction_id = ?", existing.ReactionID).Update("reaction", "dislike").Error; err != nil {
				tx.Rollback()
				return errno.NewErrNo(errno.InternalDatabaseErrorCode, "切换为点踩失败: "+err.Error())
			}
			if err := tx.Table(constants.CourseCommentTableName).Where("comment_id = ? AND likes > 0", commentID).Update("likes", gorm.Expr("likes - 1")).Error; err != nil {
				tx.Rollback()
				return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新点赞数失败: "+err.Error())
			}
		}

	case "cancel_like":
		if hasExisting && existing.Reaction == "like" {
			if err := tx.Table(constants.CourseCommentReactionTableName).Where("reaction_id = ?", existing.ReactionID).Delete(&CourseCommentReaction{}).Error; err != nil {
				tx.Rollback()
				return errno.NewErrNo(errno.InternalDatabaseErrorCode, "撤销课程评论点赞失败: "+err.Error())
			}
			if err := tx.Table(constants.CourseCommentTableName).Where("comment_id = ? AND likes > 0", commentID).Update("likes", gorm.Expr("likes - 1")).Error; err != nil {
				tx.Rollback()
				return errno.NewErrNo(errno.InternalDatabaseErrorCode, "更新点赞数失败: "+err.Error())
			}
		}

	case "cancel_dislike":
		if hasExisting && existing.Reaction == "dislike" {
			if err := tx.Table(constants.CourseCommentReactionTableName).Where("reaction_id = ?", existing.ReactionID).Delete(&CourseCommentReaction{}).Error; err != nil {
				tx.Rollback()
				return errno.NewErrNo(errno.InternalDatabaseErrorCode, "撤销课程评论点踩失败: "+err.Error())
			}
		}

	default:
		tx.Rollback()
		return errno.ParamVerifyError.WithMessage("操作类型无效")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "提交课程评论反应事务失败: "+err.Error())
	}
	return nil
}

// ReactCourseCommentAsync 异步处理课程评论反应（点赞/点踩等）
func ReactCourseCommentAsync(ctx context.Context, userID, commentID int64, action string) chan error {
	pool := GetAsyncPool()
	return pool.Submit(func() error {
		return ReactCourseComment(ctx, userID, commentID, action)
	})
}
