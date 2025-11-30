namespace go resource
include "model.thrift"

// 搜索资源请求
struct SearchResourceReq {
    1: optional string keyword,          
    2: optional i64 tagId,               
    3: optional string sortBy, 
    4: optional i64 course_id,  
    5: required i32 page_size, 
    6: required i32 page_num,             
}

struct SearchResourceResp {
    1: required model.BaseResp baseResp,
    2: required list<model.Resource> resources,
    3: required i32 total, 
}

// 上传资源请求
struct UploadResourceReq {
    1: required binary file, 
    2: required string title,  
    3: optional string description,
    4: required i64 course_id, 
    5: optional list<string> tags,
}

struct UploadResourceResp {
    1: required model.BaseResp baseResp,
    2: optional model.Resource resource,
}

// 下载资源请求
struct DownloadResourceReq {
    1: required i64 resource_id (api.path="resource_id"),
}

struct DownloadResourceResp {
    1: required model.BaseResp baseResp,
    2: required string downloadUrl,
}

// 举报资源请求
struct ReportResourceReq {
    1: required i64 resource_id (api.path="resource_id"),
    2: required string reason (api.body="reason"),
}

struct ReportResourceResp {
    1: required model.BaseResp baseResp,
}

// 获取资源信息请求
struct GetResourceReq {
    1: required i64 resource_id (api.path="resource_id"), 
}

struct GetResourceResp {
    1: required model.BaseResp baseResp,
    2: optional model.Resource resource,
}

// 提交资源评分请求
struct SubmitResourceRatingReq {
    1: required i64 resource_id (api.path="resource_id"),
    2: required double rating,
}

struct SubmitResourceRatingResp {
    1: required model.BaseResp baseResp,
}

// 删除资源评分请求
struct DeleteResourceRatingReq {
    1: required i64 rating_id (api.path="rating_id"),
}

struct DeleteResourceRatingResp {
    1: required model.BaseResp baseResp,
}

// 提交资源评价请求
struct SubmitResourceCommentReq {
    1: required i64 resource_id (api.path="resource_id"),
    2: required string content,
    3: optional i64 parentId,
}

struct SubmitResourceCommentResp {
    1: required model.BaseResp baseResp,
}

// 删除资源评价请求
struct DeleteResourceCommentReq {
    1: required i64 comment_id (api.path="comment_id"),
}

struct DeleteResourceCommentResp {
    1: required model.BaseResp baseResp,
}

// 获取资源评论列表请求
struct GetResourceCommentsReq {
    1: required i64 resource_id (api.path="resource_id"),
    2: required i32 page_size,
    3: required i32 page_num,
    4: optional string sortBy, // latest, hottest
}

struct GetResourceCommentsResp {
    1: required model.BaseResp baseResp,
    2: required list<model.ResourceCommentWithUser> comments,
    3: required i32 total,
}

struct SubmitResourceCommentReactionReq {
    1: required i64 comment_id (api.path="comment_id"),
    2: required string action (api.body="action"),
}

struct SubmitResourceCommentReactionResp {
    1: required model.BaseResp baseResp,
}

// 资源服务
service ResourceService {
    SearchResourceResp searchResources(1: SearchResourceReq req)(api.get="/api/resources/search"),
    UploadResourceResp uploadResource(1: UploadResourceReq req)(api.post="/api/resources"),
    DownloadResourceResp downloadResource(1: DownloadResourceReq req)(api.get="/api/resources/:resource_id/download"),
    ReportResourceResp reportResource(1: ReportResourceReq req)(api.post="/api/resources/:resource_id/report"),
    GetResourceResp getResource(1: GetResourceReq req)(api.get="/api/resources/:resource_id"),
    
    // 资源评分相关API
    SubmitResourceRatingResp submitResourceRating(1: SubmitResourceRatingReq req)(api.post="/api/resource_ratings/:resource_id"),
    DeleteResourceRatingResp deleteResourceRating(1: DeleteResourceRatingReq req)(api.delete="/api/resource_ratings/:rating_id"),
    
    // 资源评论相关API
    SubmitResourceCommentResp submitResourceComment(1: SubmitResourceCommentReq req)(api.post="/api/resource_comments/:resource_id"),
    DeleteResourceCommentResp deleteResourceComment(1: DeleteResourceCommentReq req)(api.delete="/api/resource_comments/:comment_id"),
    GetResourceCommentsResp getResourceComments(1: GetResourceCommentsReq req)(api.get="/api/resource/:resource_id/comment"),

    SubmitResourceCommentReactionResp reactResourceComment(1: SubmitResourceCommentReactionReq req)(api.post="/api/resource_comments/:comment_id/likes"),
}

struct AdminDeleteResourceCommentReq{
    required i64 comment_id(api.path="comment_id"),
}
struct AdminDeleteResourceCommentResp{
    required model.BaseResp base_resp,
}


struct AdminDeleteResourceRatingReq{
    required i64 rating_id(api.path="rating_id"),
}
struct AdminDeleteResourceRatingResp{
    required model.BaseResp base_resp,
}

struct AdminDeleteResourceReq{
    required i64 resource_id(api.path="resource_id"),
}
struct AdminDeleteResourceResp{
    required model.BaseResp base_resp,
}

service AdminResourceService {
    AdminDeleteResourceCommentResp AdminDeleteResourceComment(1:AdminDeleteResourceCommentReq req)(api.delete="/api/admin/resource_comments/:comment_id"),
    AdminDeleteResourceRatingResp AdminDeleteResourceRating(1:AdminDeleteResourceRatingReq req)(api.delete="/api/admin/resource_ratings/:rating_id"),
    AdminDeleteResourceResp AdminDeleteResource(1:AdminDeleteResourceReq req)(api.delete="/api/admin/resources/:resource_id"),

}
