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
  required i64 course_id  
}

struct GetCourseDetailResp {
  required model.BaseResp baseResponse;
  optional model.Course course;  
}

// 获取课程资源列表
struct GetCourseResourceListReq {
  required i64 course_id  
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
  required i64 course_id
  optional string sort_by = "latest" 
  required i32 page_size
  required i32 page_num
}

struct GetCourseCommentsResp {
  required model.BaseResp baseResponse;
  optional list<model.CourseComment> comments;
}


// 提交课程评分
struct SubmitCourseRatingReq {
  required i64 rating_id
  required i64 user_id
  required i64 course_id
  required i64 recommendation  
  required string difficulty   
  required i64 workload        
  required i64 usefulness      
  optional bool is_visible = true  
}

struct SubmitCourseRatingResp {
  required model.BaseResp baseResponse;
  optional model.CourseRating rating;  
}


// 提交课程评论
struct SubmitCourseCommentReq {
  required i64 course_id
  required string content      
  optional i64 parent_id = 0   
  optional bool is_visible = true  
}

struct SubmitCourseCommentResp {
  required model.BaseResp baseResponse;
  optional model.CourseComment comment;  
}

// 删除课程评论
struct DeleteCourseCommentReq {
  required i64 comment_id
}

struct DeleteCourseCommentResp {
  required model.BaseResp baseResponse;
}

// 删除课程评分
struct DeleteCourseRatingReq {
  required i64 rating_id
}

struct DeleteCourseRatingResp {
  required model.BaseResp baseResponse;
}

service CourseService {
  SearchResp search(1: SearchReq req)(api.get="/api/courses/search"),
  GetCourseDetailResp getCourseDetail(1: GetCourseDetailReq req)(api.get="/api/courses/{course_id}"),
  GetCourseResourceListResp getCourseResourceList(1: GetCourseResourceListReq req)(api.get="/api/courses/{course_id}/resources"),
  GetCourseCommentsResp getCourseComments(1: GetCourseCommentsReq req)(api.get="/api/courses/{course_id}/comments"),
  SubmitCourseRatingResp submitCourseRating(1: SubmitCourseRatingReq req)(api.post="/api/course_ratings/{rating_id}"),
  SubmitCourseCommentResp submitCourseComment(1: SubmitCourseCommentReq req)(api.post="/api/courses/{course_id}/comments"),
  DeleteCourseCommentResp deleteCourseComment(1: DeleteCourseCommentReq req)(api.delete="/api/courses_comments/{comment_id}"),
  DeleteCourseRatingResp deleteCourseRating(1: DeleteCourseRatingReq req)(api.delete="/api/course_ratings/{rating_id}"),
}


