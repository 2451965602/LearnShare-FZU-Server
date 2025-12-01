namespace go course
include "model.thrift"


// 搜索课程
struct SearchReq{
  required i32 page_size
  required i32 page_num
  optional string keywords    
  optional i64 college_id     
  optional string grade       
  optional double min_rating  
}

struct SearchResp {
  required model.BaseResp baseResponse;
  optional list<model.Course> courses; 
}

// 获取课程详情
struct GetCourseDetailReq {
  required i64 course_id  (api.path="course_id");
}

struct GetCourseDetailResp {
  required model.BaseResp baseResponse;
  optional model.Course course;  
}

// 获取课程资源列表
struct GetCourseResourceListReq {
  required i64 course_id  (api.path="course_id");
  required i32 page_num   
  required i32 page_size  
  optional string type     
  optional string status   
}

struct GetCourseResourceListResp {
  required model.BaseResp baseResponse;
  optional list<model.Resource> resources;
}


// 获取课程评论列表
struct GetCourseCommentsReq {
  required i64 course_id (api.path="course_id");
  optional string sort_by = "latest" 
  required i32 page_size
  required i32 page_num
}

struct GetCourseCommentsResp {
  required model.BaseResp baseResponse;
  optional list<model.CourseCommentWithUser> comments;
}


// 提交课程评分
struct SubmitCourseRatingReq {
  required i64 course_id (api.path="course_id");
  required double rating
  required i32 difficulty
  required i32 workload
  required i32 usefulness
}

struct SubmitCourseRatingResp {
  required model.BaseResp baseResponse;
  required model.CourseRating rating;  
}


// 提交课程评论
struct SubmitCourseCommentReq {
  required i64 course_id (api.path="course_id");
  required string contents      
  optional i64 parent_id
}

struct SubmitCourseCommentResp {
  required model.BaseResp baseResponse;
  required model.CourseComment comment;  
}

// 删除课程评论
struct DeleteCourseCommentReq {
  required i64 comment_id (api.path="comment_id"); 
}

struct DeleteCourseCommentResp {
  required model.BaseResp baseResponse;
}

// 删除课程评分
struct DeleteCourseRatingReq {
  required i64 rating_id (api.path="rating_id"); 
}

struct DeleteCourseRatingResp {
  required model.BaseResp baseResponse;
}

struct SubmitCourseCommentReactionReq {
    1: required i64 comment_id (api.path="comment_id"),
    2: required string action (api.body="action"),  // 例如 "like" 或 "dislike"
}

struct SubmitCourseCommentReactionResp {
    1: required model.BaseResp baseResp,
}

struct GetCourseImageReq{
    required string course_name
}

struct GetCourseImageResp{
    required  model.BaseResp baseResponse;
    required string url
}

service CourseService {
  SearchResp search(1: SearchReq req)(api.get="/api/courses/search"),
  GetCourseDetailResp getCourseDetail(1: GetCourseDetailReq req)(api.get="/api/courses/:course_id"),
  GetCourseResourceListResp getCourseResourceList(1: GetCourseResourceListReq req)(api.get="/api/courses/:course_id/resources"),
  GetCourseCommentsResp getCourseComments(1: GetCourseCommentsReq req)(api.get="/api/courses/:course_id/comments"),
  SubmitCourseRatingResp submitCourseRating(1: SubmitCourseRatingReq req)(api.post="/api/course_ratings/:course_id"),
  SubmitCourseCommentResp submitCourseComment(1: SubmitCourseCommentReq req)(api.post="/api/course_comments/:course_id"),
  DeleteCourseCommentResp deleteCourseComment(1: DeleteCourseCommentReq req)(api.delete="/api/courses_comments/:comment_id"),
  DeleteCourseRatingResp deleteCourseRating(1: DeleteCourseRatingReq req)(api.delete="/api/course_ratings/:rating_id"),

  SubmitCourseCommentReactionResp reactCourseComment(  1: SubmitCourseCommentReactionReq req)(api.post="/api/course_comments/:comment_id/likes"),

  GetCourseImageResp GetCourseImage(1:GetCourseImageReq req)(api.get="/api/course")
}

struct AdminDeleteCourseCommentReq{
    required i64 comment_id(api.path="comment_id"),
}
struct AdminDeleteCourseCommentResp{
    required model.BaseResp base_resp,
}

struct AdminDeleteCourseRatingReq{
    required i64 rating_id(api.path="rating_id"),
}
struct AdminDeleteCourseRatingResp{
    required model.BaseResp base_resp,
}

struct AdminDeleteCourseReq{
    required i64 course_id(api.path="course_id"),
}
struct AdminDeleteCourseResp{
    required model.BaseResp base_resp,
}

service AdminCourseService{
    AdminDeleteCourseCommentResp AdminDeleteCourseComment(1:AdminDeleteCourseCommentReq req)(api.delete="/api/admin/course_comments/:comment_id"),
    AdminDeleteCourseRatingResp AdminDeleteCourseRating(1:AdminDeleteCourseRatingReq req)(api.delete="/api/admin/course_ratings/:rating_id"),
    AdminDeleteCourseResp AdminDeleteCourse(1:AdminDeleteCourseReq req)(api.delete="/api/admin/courses/:course_id"),

}



