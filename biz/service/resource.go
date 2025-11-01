package service

import (
	"LearnShare/biz/dal/db"
	model "LearnShare/biz/model/module"
	"LearnShare/biz/model/resource"
	"context"
)

// SearchResourcesService 封装了搜索资源的服务
type SearchResourcesService struct {
	ctx context.Context
}

// NewSearchResourcesService 创建一个新的 SearchResourcesService
func NewSearchResourcesService(ctx context.Context) *SearchResourcesService {
	return &SearchResourcesService{ctx: ctx}
}

// SearchResources 执行资源搜索
func (s *SearchResourcesService) SearchResources(req *resource.SearchResourceReq) ([]*model.Resource, int64, error) {
	resources, total, err := db.SearchResources(s.ctx, req.Keyword, req.TagId, req.CourseId, req.SortBy, int(req.PageNum), int(req.PageSize))
	if err != nil {
		return nil, 0, err
	}

	// 将 db.Resource 转换为 model.Resource
	var modelResources []*model.Resource
	for _, r := range resources {
		var tags []*model.ResourceTag
		for _, t := range r.Tags {
			tags = append(tags, &model.ResourceTag{
				TagId:   t.TagID,
				TagName: t.TagName,
			})
		}

		modelResources = append(modelResources, &model.Resource{
			ResourceId:    r.ResourceID,
			Title:         r.Title,
			Description:   &r.Description,
			FilePath:      r.FilePath,
			FileType:      r.FileType,
			FileSize:      r.FileSize,
			UploaderId:    r.UploaderID,
			CourseId:      r.CourseID,
			DownloadCount: r.DownloadCount,
			AverageRating: r.AverageRating,
			RatingCount:   r.RatingCount,
			Status:        r.Status,
			CreatedAt:     r.CreatedAt.Unix(),
			Tags:          tags,
		})
	}

	return modelResources, total, nil
}

// GetResourceService 封装了获取单个资源信息的服务
type GetResourceService struct {
	ctx context.Context
}

// NewGetResourceService 创建一个新的 GetResourceService
func NewGetResourceService(ctx context.Context) *GetResourceService {
	return &GetResourceService{ctx: ctx}
}

// GetResource 执行获取单个资源信息
func (s *GetResourceService) GetResource(req *resource.GetResourceReq) (*model.Resource, error) {
	resource, err := db.GetResourceByID(s.ctx, req.ResourceId)
	if err != nil {
		return nil, err
	}

	// 将 db.Resource 转换为 model.Resource
	var tags []*model.ResourceTag
	for _, t := range resource.Tags {
		tags = append(tags, &model.ResourceTag{
			TagId:   t.TagID,
			TagName: t.TagName,
		})
	}

	return &model.Resource{
		ResourceId:    resource.ResourceID,
		Title:         resource.Title,
		Description:   &resource.Description,
		FilePath:      resource.FilePath,
		FileType:      resource.FileType,
		FileSize:      resource.FileSize,
		UploaderId:    resource.UploaderID,
		CourseId:      resource.CourseID,
		DownloadCount: resource.DownloadCount,
		AverageRating: resource.AverageRating,
		RatingCount:   resource.RatingCount,
		Status:        resource.Status,
		CreatedAt:     resource.CreatedAt.Unix(),
		Tags:          tags,
	}, nil
}

// GetResourceCommentsService 封装了获取资源评论列表的服务
type GetResourceCommentsService struct {
	ctx context.Context
}

// NewGetResourceCommentsService 创建一个新的 GetResourceCommentsService
func NewGetResourceCommentsService(ctx context.Context) *GetResourceCommentsService {
	return &GetResourceCommentsService{ctx: ctx}
}

// GetResourceComments 执行获取资源评论列表
func (s *GetResourceCommentsService) GetResourceComments(req *resource.GetResourceCommentsReq) ([]*model.ResourceComment, int64, error) {
	// 调用数据库层获取评论数据
	comments, total, err := db.GetResourceComments(s.ctx, req.ResourceId, req.SortBy, int(req.PageNum), int(req.PageSize))
	if err != nil {
		return nil, 0, err
	}

	// 将 db.ResourceComment 转换为 model.ResourceComment
	var modelComments []*model.ResourceComment
	for _, comment := range comments {
		modelComments = append(modelComments, &model.ResourceComment{
			CommentId:  comment.CommentID,
			UserId:     comment.UserID,
			ResourceId: comment.ResourceID,
			Content:    comment.Content,
			ParentId:   func() int64 { if comment.ParentID != nil { return *comment.ParentID } else { return 0 } }(),
			Likes:      comment.Likes,
			IsVisible:  comment.IsVisible,
			Status:     func() model.ResourceCommentStatus { 
			status, _ := model.ResourceCommentStatusFromString(comment.Status)
			return status
		}(),
			CreatedAt:  comment.CreatedAt.Unix(),
		})
	}

	return modelComments, total, nil
}

// SubmitResourceRatingService 封装了提交资源评分的服务
type SubmitResourceRatingService struct {
	ctx context.Context
}

// NewSubmitResourceRatingService 创建一个新的 SubmitResourceRatingService
func NewSubmitResourceRatingService(ctx context.Context) *SubmitResourceRatingService {
	return &SubmitResourceRatingService{ctx: ctx}
}

// SubmitResourceRating 执行提交资源评分
func (s *SubmitResourceRatingService) SubmitResourceRating(req *resource.SubmitResourceRatingReq, userID int64) (*model.ResourceRating, error) {
	// 调用数据库层提交评分，使用rating字段
	rating, err := db.SubmitResourceRating(s.ctx, userID, req.ResourceId, req.Rating)
	if err != nil {
		return nil, err
	}

	// 将 db.ResourceRating 转换为 model.ResourceRating
	return &model.ResourceRating{
		RatingId:       rating.RatingID,
		UserId:         rating.UserID,
		ResourceId:     rating.ResourceID,
		Recommendation: rating.Recommendation * 10, // 转换为0-50的浮点数
		IsVisible:      rating.IsVisible,
		CreatedAt:      rating.CreatedAt.Unix(),
	}, nil
}

// SubmitResourceCommentService 封装了提交资源评论的服务
type SubmitResourceCommentService struct {
	ctx context.Context
}

// NewSubmitResourceCommentService 创建一个新的 SubmitResourceCommentService
func NewSubmitResourceCommentService(ctx context.Context) *SubmitResourceCommentService {
	return &SubmitResourceCommentService{ctx: ctx}
}

// SubmitResourceComment 执行提交资源评论
func (s *SubmitResourceCommentService) SubmitResourceComment(req *resource.SubmitResourceCommentReq, userID int64) (*model.ResourceComment, error) {
	// 处理父评论ID
	var parentID *int64
	if req.IsSetParentId() && req.ParentId != nil && *req.ParentId != 0 {
		parentID = req.ParentId
	}

	// 调用数据库层提交评论
	comment, err := db.SubmitResourceComment(s.ctx, userID, req.ResourceId, req.Content, parentID)
	if err != nil {
		return nil, err
	}

	// 将 db.ResourceComment 转换为 model.ResourceComment
	return &model.ResourceComment{
		CommentId:  comment.CommentID,
		UserId:     comment.UserID,
		ResourceId: comment.ResourceID,
		Content:    comment.Content,
		ParentId:   func() int64 { if comment.ParentID != nil { return *comment.ParentID } else { return 0 } }(),
		Likes:      comment.Likes,
		IsVisible:  comment.IsVisible,
		Status:     func() model.ResourceCommentStatus { 
			status, _ := model.ResourceCommentStatusFromString(comment.Status)
			return status
		}(),
		CreatedAt:  comment.CreatedAt.Unix(),
	}, nil
}
