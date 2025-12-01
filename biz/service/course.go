package service

import (
	"LearnShare/biz/dal/db"
	"LearnShare/biz/model/course"
	"LearnShare/biz/model/module"
	"LearnShare/pkg/errno"
	"bytes"
	"context"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/fogleman/gg"
)

type CourseService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewCourseService(ctx context.Context, c *app.RequestContext) *CourseService {
	return &CourseService{ctx: ctx, c: c}
}

func (s *CourseService) Search(req *course.SearchReq) ([]*module.Course, error) {

	// 使用strings包处理指针类型的参数
	keywords := ""
	if req.Keywords != nil {
		keywords = strings.TrimSpace(*req.Keywords)
	}

	grade := ""
	if req.Grade != nil {
		grade = strings.TrimSpace(*req.Grade)
	}

	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 调用数据库查询课程
	courses, err := db.SearchCourses(s.ctx, keywords, grade, int(req.PageNum), int(req.PageSize))
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "搜索课程失败: "+err.Error())
	}

	// 转换为module.Course列表
	var courseModules []*module.Course
	for _, c := range courses {
		courseModules = append(courseModules, c.ToCourseModule())
	}

	return courseModules, nil
}

func (s *CourseService) GetCourseDetail(req *course.GetCourseDetailReq) (*module.Course, error) {
	// 获取课程详情
	courseDetail, err := db.GetCourseByID(s.ctx, req.CourseID)
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "获取课程详情失败: "+err.Error())
	}

	return courseDetail.ToCourseModule(), nil
}

func (s *CourseService) GetCourseResourceList(req *course.GetCourseResourceListReq) ([]*module.Resource, error) {
	// 处理指针类型的参数
	var resourceType string
	if req.Type != nil {
		resourceType = *req.Type
	}

	var status string
	if req.Status != nil {
		status = *req.Status
	}
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 获取课程资源列表 - 使用正确的字段名
	resources, err := db.GetCourseResources(s.ctx, req.CourseID, resourceType, status, int(req.PageNum), int(req.PageSize)) // 改为 CourseID
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "获取课程资源失败: "+err.Error())
	}

	// 转换为module.Resource列表
	var resourceModules []*module.Resource
	for _, r := range resources {
		resourceModules = append(resourceModules, r.ToResourceModule())
	}

	return resourceModules, nil
}

func (s *CourseService) GetCourseComments(req *course.GetCourseCommentsReq) ([]*module.CourseCommentWithUser, error) {
	// SortBy 是普通 string 类型，不是指针
	sortBy := req.SortBy
	if sortBy == "" {
		sortBy = "latest" // 使用默认值
	}

	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 获取课程评论列表
	comments, err := db.GetCourseCommentsByCourseID(s.ctx, req.CourseID, sortBy, int(req.PageNum), int(req.PageSize))
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "获取课程评论失败: "+err.Error())
	}

	// 转换为module.CourseComment列表
	var commentModules []*module.CourseCommentWithUser
	for _, c := range comments {
		commentModules = append(commentModules, c.ToCourseCommentWithUserModule())
	}

	return commentModules, nil
}

func (s *CourseService) SubmitCourseRating(req *course.SubmitCourseRatingReq) (*module.CourseRating, error) {
	userID := GetUidFormContext(s.c)

	if req.Rating < 0 || req.Rating > 5 {
		return nil, errno.ValidationRatingRangeInvalidError
	}

	if req.Difficulty < 1 || req.Difficulty > 5 {
		return nil, errno.NewErrNo(errno.ServiceInvalidParameter, "难度必须在1-5之间")
	}
	if req.Workload < 1 || req.Workload > 5 {
		return nil, errno.NewErrNo(errno.ServiceInvalidParameter, "考核压力必须在1-5之间")
	}
	if req.Usefulness < 1 || req.Usefulness > 5 {
		return nil, errno.NewErrNo(errno.ServiceInvalidParameter, "实用性必须在1-5之间")
	}

	rating := &db.CourseRating{
		UserID:         userID,
		CourseID:       req.CourseID,
		Recommendation: req.Rating,
		Difficulty:     uint8(req.Difficulty),
		Workload:       uint8(req.Workload),
		Usefulness:     uint8(req.Usefulness),
		IsVisible:      true,
	}

	saved, err := db.SubmitCourseRating(s.ctx, rating)
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "提交课程评分失败: "+err.Error())
	}

	return saved.ToCourseRatingModule(), nil
}

func (s *CourseService) SubmitCourseComment(req *course.SubmitCourseCommentReq) (*module.CourseComment, error) {
	// 获取用户ID
	userID := GetUidFormContext(s.c)

	comment := &db.CourseComment{
		CourseID: req.CourseID,
		UserID:   userID,
		Content:  req.Contents,
		ParentID: req.ParentID,
		Status:   "normal",
	}

	// 使用异步提交评论
	if err := db.SubmitCourseComment(s.ctx, comment); err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "提交评论失败: "+err.Error())
	}

	return comment.ToCourseCommentModule(), nil
}

func (s *CourseService) DeleteCourseComment(req *course.DeleteCourseCommentReq) error {
	// 使用异步删除评论
	errChan := db.DeleteCourseCommentAsync(s.ctx, req.CommentID)
	if err := <-errChan; err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "删除评论失败: "+err.Error())
	}

	return nil
}

func (s *CourseService) DeleteCourseRating(req *course.DeleteCourseRatingReq) error {
	// 使用异步删除评分
	errChan := db.DeleteCourseRatingAsync(s.ctx, req.RatingID)
	if err := <-errChan; err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "删除评分失败: "+err.Error())
	}

	return nil
}

func (s *CourseService) ReactCourseComment(commentID int64, action string) error {
	if commentID <= 0 {
		return errno.ParamVerifyError
	}
	switch action {
	case "like", "dislike", "cancel_like", "cancel_dislike":
	default:
		return errno.ParamVerifyError.WithMessage("操作类型无效")
	}

	userID := GetUidFormContext(s.c)

	errChan := db.ReactCourseCommentAsync(s.ctx, userID, commentID, action)
	if err := <-errChan; err != nil {
		return err
	}
	return nil
}

// AdminDeleteCourse 管理员硬删除课程（包括关联资源、评论、评分等）
func (s *CourseService) AdminDeleteCourse(req *course.AdminDeleteCourseReq) error {
	if req.CourseID <= 0 {
		return errno.NewErrNo(errno.ServiceInvalidParameter, "课程ID无效")
	}

	// 调用数据库层执行硬删除（需确保 db.AdminDeleteCourse 已实现级联删除或事务清理）
	if err := db.AdminDeleteCourse(s.ctx, req.CourseID); err != nil {
		return err
	}

	return nil
}

// AdminDeleteCourseComment 管理员硬删除课程评论
func (s *CourseService) AdminDeleteCourseComment(req *course.AdminDeleteCourseCommentReq) error {
	if req.CommentID <= 0 {
		return errno.NewErrNo(errno.ServiceInvalidParameter, "评论ID无效")
	}

	if err := db.AdminDeleteCourseComment(s.ctx, req.CommentID); err != nil {
		return err
	}
	return nil
}

// AdminDeleteCourseRating 管理员硬删除课程评分
func (s *CourseService) AdminDeleteCourseRating(req *course.AdminDeleteCourseRatingReq) error {
	if req.RatingID <= 0 {
		return errno.NewErrNo(errno.ServiceInvalidParameter, "评分ID无效")
	}

	if err := db.AdminDeleteCourseRating(s.ctx, req.RatingID); err != nil {
		return err
	}
	return nil
}

func (s *CourseService) GetCourseImage(name string) (string, error) {
	const width, height = 800, 400
	dc := gg.NewContext(width, height)

	// 创建渐变背景
	grad := gg.NewLinearGradient(0, 0, width, height)
	grad.AddColorStop(0, color.RGBA{R: 99, G: 102, B: 241, A: 255})
	grad.AddColorStop(1, color.RGBA{R: 168, G: 85, B: 247, A: 255})
	dc.SetFillStyle(grad)
	dc.DrawRectangle(0, 0, float64(width), float64(height))
	dc.Fill()

	// 绘制文字阴影

	err := dc.LoadFontFace("config/font/msyh.ttc", 80)
	if err != nil {
		return "", err
	}

	dc.SetColor(color.RGBA{A: 100})
	dc.DrawStringAnchored(name, width/2+4, height/2+4, 0.5, 0.5)

	// 绘制主文字
	dc.SetColor(color.White)
	dc.DrawStringAnchored(name, width/2, height/2, 0.5, 0.5)

	// 编码为base64
	var buf bytes.Buffer
	if err := png.Encode(&buf, dc.Image().(image.Image)); err != nil {
		return "", err
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
