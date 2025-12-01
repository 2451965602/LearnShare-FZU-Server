package db

import (
	"LearnShare/biz/model/module"
	"fmt"
	"time"
)

type User struct {
	UserID          int64     `json:"user_id" db:"user_id"`
	Username        string    `json:"username" db:"username"`
	PasswordHash    string    `json:"-" db:"password_hash"`
	Email           string    `json:"email" db:"email"`
	CollegeID       *int64    `json:"college_id,omitempty" db:"college_id"`
	MajorID         *int64    `json:"major_id,omitempty" db:"major_id"`
	AvatarURL       *string   `json:"avatar_url,omitempty" db:"avatar_url"`
	ReputationScore int64     `json:"reputation_score" db:"reputation_score"`
	RoleID          int64     `json:"role_id" db:"role_id"`
	Status          string    `json:"status" db:"status"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

func (u User) ToUserModule() *module.User {
	user := &module.User{
		UserId:          u.UserID,
		Username:        u.Username,
		Email:           u.Email,
		ReputationScore: u.ReputationScore,
		RoleId:          u.RoleID,
		Status:          u.Status,
		CreatedAt:       u.CreatedAt.Unix(),
		UpdatedAt:       u.UpdatedAt.Unix(),
	}

	if u.AvatarURL != nil {
		user.AvatarURL = *u.AvatarURL
	}
	if u.CollegeID != nil {
		user.CollegeID = *u.CollegeID
	}
	if u.MajorID != nil {
		user.MajorID = *u.MajorID
	}
	return user
}

// Course 相关结构体
type Course struct {
	CourseID    int64     `json:"course_id" db:"course_id"`
	CourseName  string    `json:"course_name" db:"course_name"`
	TeacherID   int64     `json:"teacher_id" db:"teacher_id"`
	Credit      float64   `json:"credit" db:"credit"`
	MajorID     int64     `json:"major_id" db:"major_id"`
	Grade       string    `json:"grade" db:"grade"`
	Description *string   `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func (c Course) ToCourseModule() *module.Course {

	var description string
	if c.Description != nil {
		description = *c.Description
	} else {
		description = ""
	}

	course := &module.Course{
		CourseId:    c.CourseID,
		CourseName:  c.CourseName,
		TeacherId:   c.TeacherID,
		Credit:      c.Credit,
		MajorId:     c.MajorID,
		Grade:       c.Grade,
		Description: description,
		CreatedAt:   c.CreatedAt.Unix(),
		UpdatedAt:   c.UpdatedAt.Unix(),
	}
	return course
}

type Resource struct {
	ResourceID    int64         `gorm:"primaryKey;autoIncrement"`
	ResourceName  string        `gorm:"column:resource_name;size:255;not null"`
	Description   string        `gorm:"type:text"`
	FilePath      string        `gorm:"column:resource_url;size:255;not null"`
	FileType      string        `gorm:"column:type;size:50;not null"`
	FileSize      int64         `gorm:"column:size;not null"`
	UploaderID    int64         `gorm:"not null"`
	CourseID      int64         `gorm:"not null"`
	DownloadCount int64         `gorm:"default:0"`
	AverageRating float64       `gorm:"default:0.0"`
	RatingCount   int64         `gorm:"default:0"`
	Status        string        `gorm:"type:enum('normal','low_quality','pending_review');default:'pending_review'"`
	CreatedAt     time.Time     `gorm:"autoCreateTime"`
	Tags          []ResourceTag `gorm:"many2many:resource_tags;joinForeignKey:resource_id;joinReferences:tag_id"`
}

// ToResourceModule 将db.Resource转换为model.Resource
func (r Resource) ToResourceModule() *module.Resource {
	var tags []*module.ResourceTag
	seen := make(map[int64]struct{})
	for _, t := range r.Tags {
		if _, ok := seen[t.TagID]; ok {
			continue
		}
		seen[t.TagID] = struct{}{}
		tags = append(tags, t.ToResourceTagModule())
	}

	if tags == nil {
		tags = make([]*module.ResourceTag, 0)
	}
	return &module.Resource{
		ResourceId:    r.ResourceID,
		Title:         r.ResourceName,
		Description:   &r.Description,
		FileType:      r.FileType,
		FileSize:      r.FileSize,
		UploaderId:    r.UploaderID,
		CourseId:      r.CourseID,
		DownloadCount: r.DownloadCount,
		AverageRating: r.AverageRating,
		RatingCount:   r.RatingCount,
		Status:        convertStatus(r.Status),
		CreatedAt:     r.CreatedAt.Unix(),
		Tags:          tags,
	}
}

func convertStatus(status string) int32 {
	// 与接口文档保持一致：0 待审核, 1 已发布, 2 已拒绝
	switch status {
	case "pending_review":
		return 0
	case "normal":
		return 1
	case "low_quality":
		return 2
	default:
		return 0
	}
}

type ResourceTag struct {
	TagID   int64  `gorm:"primaryKey;column:tag_id"`
	TagName string `gorm:"column:tag_name;size:50;unique;not null"`
}

// TableName 指定标签表名为 tags
func (ResourceTag) TableName() string { return "tags" }

// ToResourceTagModule 将db.ResourceTag转换为model.ResourceTag
func (t ResourceTag) ToResourceTagModule() *module.ResourceTag {
	return &module.ResourceTag{
		TagId:   t.TagID,
		TagName: t.TagName,
	}
}

type ResourceCommentrow struct {
	CommentID  int64     `gorm:"column:comment_id"`
	UserID     int64     `gorm:"column:user_id"`
	ResourceID int64     `gorm:"column:resource_id"`
	Content    string    `gorm:"column:content"`
	ParentID   *int64    `gorm:"column:parent_id"`
	Likes      int64     `gorm:"column:likes"`
	IsVisible  bool      `gorm:"column:is_visible"`
	Status     string    `gorm:"column:status"`
	CreatedAt  time.Time `gorm:"column:created_at"`

	UUserID          *int64     `gorm:"column:u_user_id"`
	UUsername        *string    `gorm:"column:u_username"`
	UPasswordHash    *string    `gorm:"column:u_password_hash"`
	UEmail           *string    `gorm:"column:u_email"`
	UCollegeID       *int64     `gorm:"column:u_college_id"`
	UMajorID         *int64     `gorm:"column:u_major_id"`
	UAvatarURL       *string    `gorm:"column:u_avatar_url"`
	UReputationScore *int64     `gorm:"column:u_reputation_score"`
	URoleID          *int64     `gorm:"column:u_role_id"`
	UStatus          *string    `gorm:"column:u_status"`
	UCreatedAt       *time.Time `gorm:"column:u_created_at"`
	UUpdatedAt       *time.Time `gorm:"column:u_updated_at"`
}

// CourseRating 课程评分
type CourseRating struct {
	RatingID       int64     `json:"rating_id" db:"rating_id"`
	UserID         int64     `json:"user_id" db:"user_id"`
	CourseID       int64     `json:"course_id" db:"course_id"`
	Recommendation float64   `json:"recommendation" db:"recommendation"`
	Difficulty     uint8     `json:"difficulty" db:"difficulty"`
	Workload       uint8     `json:"workload" db:"workload"`
	Usefulness     uint8     `json:"usefulness" db:"usefulness"`
	IsVisible      bool      `json:"is_visible" db:"is_visible"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

func (r CourseRating) ToCourseRatingModule() *module.CourseRating {
	return &module.CourseRating{
		RatingId:       r.RatingID,
		UserId:         r.UserID,
		CourseId:       r.CourseID,
		Recommendation: int32(r.Recommendation),
		Difficulty:     fmt.Sprintf("%d", r.Difficulty),
		Workload:       int32(r.Workload),
		Usefulness:     int32(r.Usefulness),
		IsVisible:      r.IsVisible,
		CreatedAt:      r.CreatedAt.Unix(),
	}
}

type CourseComment struct {
	CommentID int64     `json:"comment_id" db:"comment_id"`
	CourseID  int64     `json:"course_id" db:"course_id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Content   string    `json:"content" db:"content"`
	Likes     int64     `json:"likes" db:"likes"`
	ParentID  *int64    `json:"parent_id,omitempty" db:"parent_id"` // 允许 NULL
	IsVisible bool      `json:"is_visible" db:"is_visible"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (c CourseComment) ToCourseCommentModule() *module.CourseComment {
	return &module.CourseComment{
		CommentId: c.CommentID,
		UserId:    c.UserID,
		CourseId:  c.CourseID,
		Content:   c.Content,
		ParentId:  c.ParentID, // *int64 → optional i64 ✅
		Likes:     c.Likes,
		IsVisible: c.IsVisible,
		Status:    c.Status,
		CreatedAt: c.CreatedAt.Unix(),
	}
}

// CourseCommentReaction 课程评论反应模型
type CourseCommentReaction struct {
	ReactionID int64     `gorm:"primaryKey;autoIncrement"`
	UserID     int64     `gorm:"not null"`
	CommentID  int64     `gorm:"not null"`
	Reaction   string    `gorm:"type:enum('like','dislike');not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

type CommentUserRow struct {
	CommentID int64     `gorm:"column:comment_id"`
	CourseID  int64     `gorm:"column:course_id"`
	Content   string    `gorm:"column:content"`
	Likes     int64     `json:"likes" db:"likes"`
	Status    string    `json:"status" db:"status"`
	ParentID  *int64    `gorm:"column:parent_id"`
	IsVisible bool      `gorm:"column:is_visible"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`

	UserID          *int64  `gorm:"column:u_user_id"`
	Username        *string `gorm:"column:u_username"`
	Email           *string `gorm:"column:u_email"`
	CollegeID       *int64  `gorm:"column:u_college_id"`
	MajorID         *int64  `gorm:"column:u_major_id"`
	AvatarURL       *string `gorm:"column:u_avatar_url"`
	ReputationScore *int64  `gorm:"column:u_reputation_score"`
	RoleID          *int64  `gorm:"column:u_role_id"`
	UserStatus      *string `gorm:"column:u_status"`
}

type CourseCommentWithuser struct {
	CommentID int64     `json:"comment_id" db:"comment_id"`
	CourseID  int64     `json:"course_id" db:"course_id"`
	User      User      `json:"user" db:"-"`
	Likes     int64     `json:"likes" db:"likes"`
	Content   string    `json:"content" db:"content"`
	ParentID  *int64    `json:"parent_id,omitempty" db:"parent_id"` // ← 关键修改！
	IsVisible bool      `json:"is_visible" db:"is_visible"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (c CourseCommentWithuser) ToCourseCommentWithUserModule() *module.CourseCommentWithUser {

	return &module.CourseCommentWithUser{
		CommentId: c.CommentID,
		User:      c.User.ToUserModule(),
		CourseId:  c.CourseID,
		Content:   c.Content,
		ParentId:  c.ParentID,
		Likes:     c.Likes, // 必须：Thrift required字段
		IsVisible: c.IsVisible,
		Status:    c.Status,           // 必须：Thrift required字段
		CreatedAt: c.CreatedAt.Unix(), // 必须：时间戳转换
	}
}

type ResourceTagMapping struct {
	ResourceID int64 `gorm:"primaryKey;table:resource_tags"`
	TagID      int64 `gorm:"primaryKey"`
}

// ResourceComment 资源评论模型
type ResourceComment struct {
	CommentID  int64     `gorm:"primaryKey;autoIncrement"`
	UserID     int64     `gorm:"not null"`
	ResourceID int64     `gorm:"not null"`
	Content    string    `gorm:"type:text;not null"`
	ParentID   *int64    `gorm:"default:NULL"`
	Likes      int64     `gorm:"default:0"`
	IsVisible  bool      `gorm:"default:true"`
	Status     string    `gorm:"type:enum('normal','deleted_by_user','deleted_by_admin');default:'normal'"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`

	// 关联用户信息
	User User `gorm:"foreignKey:UserID;references:UserID"`
}

// ResourceRating 资源评分模型
type ResourceRating struct {
	RatingID       int64     `gorm:"primaryKey;autoIncrement"`
	UserID         int64     `gorm:"not null"`
	ResourceID     int64     `gorm:"not null"`
	Recommendation float64   `gorm:"type:decimal(2,1);not null"`
	IsVisible      bool      `gorm:"default:true"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`

	// 关联用户信息
	User User `gorm:"foreignKey:UserID;references:UserID"`
	// 关联资源信息
	Resource Resource `gorm:"foreignKey:ResourceID;references:ResourceID"`
}

type ResourceCommentReaction struct {
	ReactionID int64     `gorm:"primaryKey;autoIncrement"`
	UserID     int64     `gorm:"not null"`
	CommentID  int64     `gorm:"not null"`
	Reaction   string    `gorm:"type:enum('like','dislike');not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

// ToResourceCommentModule 将db.ResourceComment转换为model.ResourceComment
func (c ResourceComment) ToResourceCommentModule() *module.ResourceComment {
	var parentId int64
	if c.ParentID != nil {
		parentId = *c.ParentID
	}

	var status module.ResourceCommentStatus
	status, _ = module.ResourceCommentStatusFromString(c.Status)

	return &module.ResourceComment{
		CommentId:  c.CommentID,
		UserId:     c.UserID,
		ResourceId: c.ResourceID,
		Content:    c.Content,
		ParentId:   parentId,
		Likes:      c.Likes,
		IsVisible:  c.IsVisible,
		Status:     status,
		CreatedAt:  c.CreatedAt.Unix(),
	}
}

type ResourceCommentWithUser struct {
	CommentID  int64     `gorm:"primaryKey;autoIncrement"`
	UserID     int64     `gorm:"not null"`
	ResourceID int64     `gorm:"not null"`
	Content    string    `gorm:"type:text;not null"`
	ParentID   *int64    `gorm:"default:NULL"`
	Likes      int64     `gorm:"default:0"`
	IsVisible  bool      `gorm:"default:true"`
	Status     string    `gorm:"type:enum('normal','deleted_by_user','deleted_by_admin');default:'normal'"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`

	// 关联用户信息
	User User `gorm:"foreignKey:UserID;references:UserID"`
}

func (c ResourceCommentWithUser) ToResourceCommentWithUserModule() *module.ResourceCommentWithUser {
	var parentId int64
	if c.ParentID != nil {
		parentId = *c.ParentID
	}

	var status module.ResourceCommentStatus
	status, _ = module.ResourceCommentStatusFromString(c.Status)

	return &module.ResourceCommentWithUser{
		CommentId:  c.CommentID,
		User:       c.User.ToUserModule(),
		ResourceId: c.ResourceID,
		Content:    c.Content,
		ParentId:   parentId,
		Likes:      c.Likes,
		IsVisible:  c.IsVisible,
		Status:     status,
		CreatedAt:  c.CreatedAt.Unix(),
	}
}

// ToResourceRatingModule 将db.ResourceRating转换为model.ResourceRating
func (r ResourceRating) ToResourceRatingModule() *module.ResourceRating {
	return &module.ResourceRating{
		RatingId:       r.RatingID,
		UserId:         r.UserID,
		ResourceId:     r.ResourceID,
		Recommendation: r.Recommendation * 10, // 转换为0-50的浮点数
		IsVisible:      r.IsVisible,
		CreatedAt:      r.CreatedAt.Unix(),
	}
}

type Review struct {
	ReviewID   int64      `gorm:"primaryKey;autoIncrement;column:review_id"`
	TargetID   int64      `gorm:"not null;column:target_id"`
	TargetType string     `gorm:"size:50;not null;column:target_type"`
	Reason     string     `gorm:"type:text;not null;column:reason"`
	Status     string     `gorm:"type:enum('pending','approved','rejected');default:'pending';column:status"`
	Priority   int        `gorm:"default:3;column:priority"`
	ReporterID int64      `gorm:"column:reporter_id"`
	ReviewerID *int64     `gorm:"column:reviewer_id"`
	ReviewedAt *time.Time `gorm:"column:reviewed_at"`
	CreatedAt  time.Time  `gorm:"autoCreateTime;column:created_at"`
}

// ToReviewModule 将db.Review转换为model.Review
func (r Review) ToReviewModule() *module.Review {
	var reviewerId int64
	if r.ReviewedAt != nil && r.ReviewerID != nil {
		reviewerId = *r.ReviewerID
	} else {
		reviewerId = 0
	}

	reporterId := r.ReporterID

	return &module.Review{
		ReviewId:   r.ReviewID,
		ReviewerId: reviewerId,
		ReporterId: reporterId,
		TargetId:   r.TargetID,
		TargetType: r.TargetType,
		Reason:     r.Reason,
		Status:     r.Status,
		Priority:   int64(r.Priority),
		CreatedAt:  r.CreatedAt.Unix(),
	}
}

// Permission 权限表结构
type Permission struct {
	PermissionID   int64  `json:"permission_id" db:"permission_id"`
	PermissionName string `json:"permission_name" db:"permission_name"`
	Description    string `json:"description" db:"description"`
}

func (p Permission) ToPermissionModule() *module.Permission {
	return &module.Permission{
		PermissionId:   p.PermissionID,
		PermissionName: p.PermissionName,
		Description:    p.Description,
	}
}

// Role 角色表结构
type Role struct {
	RoleID      int64  `json:"role_id" db:"role_id"`
	RoleName    string `json:"role_name" db:"role_name"`
	Description string `json:"description" db:"description"`
}

func (r Role) ToRoleModule() *module.Role {
	return &module.Role{
		RoleId:      r.RoleID,
		RoleName:    r.RoleName,
		Description: r.Description,
	}
}

// RolePermission 角色权限关联表结构
type RolePermission struct {
	RoleID       int64 `json:"role_id" db:"role_id"`
	PermissionID int64 `json:"permission_id" db:"permission_id"`
}

type RolePermRow struct {
	RoleID                int64   `gorm:"column:role_id"`
	PermissionID          *int64  `gorm:"column:permission_id"`
	PermissionName        *string `gorm:"column:permission_name"`
	PermissionDescription *string `gorm:"column:permission_description"`
}

type RoleWithPermissions struct {
	RoleID      int64        `json:"role_id" db:"role_id"`
	RoleName    string       `json:"role_name" db:"role_name"`
	Description string       `json:"description" db:"description"`
	Permissions []Permission `json:"permissions" db:"-"`
}

func (r RoleWithPermissions) ToRoleModule() *module.Role {
	role := &module.Role{
		RoleId:      r.RoleID,
		RoleName:    r.RoleName,
		Description: r.Description,
	}
	for _, p := range r.Permissions {
		role.Permissions = append(role.Permissions, p.ToPermissionModule())
	}

	return role
}

// College 学院模型
type College struct {
	CollegeID   int64  `json:"college_id" db:"college_id"`
	CollegeName string `json:"college_name" db:"college_name"`
	School      string `json:"school" db:"school"`
}

func (c College) ToCollegeModule() *module.College {
	return &module.College{
		CollegeId:   c.CollegeID,
		CollegeName: c.CollegeName,
	}
}

// Major 专业模型
type Major struct {
	MajorID   int64  `json:"major_id" db:"major_id"`
	MajorName string `json:"major_name" db:"major_name"`
	CollegeID int64  `json:"college_id" db:"college_id"`
}

func (m Major) ToMajorModule() *module.Major {
	return &module.Major{
		MajorId:   m.MajorID,
		MajorName: m.MajorName,
		CollegeId: m.CollegeID,
	}
}

// Teacher 教师模型
type Teacher struct {
	TeacherID    int64     `json:"teacher_id" db:"teacher_id"`
	Name         string    `json:"name" db:"name"`
	CollegeID    *int64    `json:"college_id" db:"college_id"`
	Introduction *string   `json:"introduction" db:"introduction"`
	Email        *string   `json:"email" db:"email"`
	AvatarURL    *string   `json:"avatar_url" db:"avatar_url"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

func (t Teacher) ToTeacherModule() *module.Teacher {
	teacher := &module.Teacher{
		TeacherId:   t.TeacherID,
		TeacherName: t.Name,
	}
	if t.CollegeID != nil {
		teacher.CollegeId = *t.CollegeID
	}
	if t.Introduction != nil {
		teacher.Introduction = *t.Introduction
	}
	if t.Email != nil {
		teacher.Email = *t.Email
	}
	if t.AvatarURL != nil {
		teacher.AvatarURL = *t.AvatarURL
	}
	return teacher
}

// Favorite 收藏模型
type Favorite struct {
	FavoriteID int64     `json:"favorite_id" db:"favorite_id"`
	UserID     int64     `json:"user_id" db:"user_id"`
	TargetID   int64     `json:"target_id" db:"target_id"`
	TargetType string    `json:"target_type" db:"target_type"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

func (f Favorite) ToFavoriteModule() *module.Favorite {
	return &module.Favorite{
		FavoriteId: f.FavoriteID,
		TargetId:   f.TargetID,
		TargetType: f.TargetType,
		CreatedAt:  f.CreatedAt.Unix(),
	}
}
