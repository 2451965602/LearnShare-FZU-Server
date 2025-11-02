namespace go resource
include "model.thrift"

// 搜索资源请求
struct SearchResourceReq {
    1: optional string keyword,          
    2: optional i64 tagId,               
    3: optional string sortBy, 
    4: optional i64 courseId,  
    5: required i32 pageSize, 
    6: required i32 pageNum,             
}

struct SearchResourceResp {
    1: required model.BaseResp baseResp,
    2: required list<model.Resource> resources,
    3: required i32 total, 
}

// 上传资源请求
struct UploadResourceReq {
    1: required binary fileData, 
    2: required string title,  
    3: optional string description,
    4: required i64 courseId, 
    5: optional list<string> tags,
}

struct UploadResourceResp {
    1: required model.BaseResp baseResp,
    2: optional model.Resource resource,
}

// 下载资源请求
struct DownloadResourceReq {
    1: required i64 resourceId, 
}

struct DownloadResourceResp {
    1: required model.BaseResp baseResp,
    2: required string downloadUrl,
}

// 举报资源请求
struct ReportResourceReq {
    1: required i64 resourceId,
}

struct ReportResourceResp {
    1: required model.BaseResp baseResp,
}

// 获取资源信息请求
struct GetResourceReq {
    1: required i64 resourceId, 
}

struct GetResourceResp {
    1: required model.BaseResp baseResp,
    2: optional model.Resource resource,
}

// 提交资源评分请求
struct SubmitResourceRatingReq {
    1: required i64 resourceId,
    2: required i64 recommendation,
}

struct SubmitResourceRatingResp {
    1: required model.BaseResp baseResp,
}

// 删除资源评分请求
struct DeleteResourceRatingReq {
    1: required i64 ratingId,
}

struct DeleteResourceRatingResp {
    1: required model.BaseResp baseResp,
}

// 提交资源评价请求
struct SubmitResourceCommentReq {
    1: required i64 resourceId,
    2: required string content,
    3: optional i64 parentId,
}

struct SubmitResourceCommentResp {
    1: required model.BaseResp baseResp,
}

// 删除资源评价请求
struct DeleteResourceCommentReq {
    1: required i64 commentId,
}

struct DeleteResourceCommentResp {
    1: required model.BaseResp baseResp,
}

// 获取资源评论列表请求
struct GetResourceCommentsReq {
    1: required i64 resourceId,
    2: required i32 pageSize,
    3: required i32 pageNum,
    4: optional string sortBy, // latest, hottest
}

struct GetResourceCommentsResp {
    1: required model.BaseResp baseResp,
    2: required list<model.ResourceComment> comments,
    3: required i32 total,
}

// 资源服务
service ResourceService {
    SearchResourceResp searchResources(1: SearchResourceReq req)(api.get="/api/resources/search"),
    UploadResourceResp uploadResource(1: UploadResourceReq req)(api.post="/api/resources"),
    DownloadResourceResp downloadResource(1: DownloadResourceReq req)(api.get="/api/resources/{resource_id}/download"),
    ReportResourceResp reportResource(1: ReportResourceReq req)(api.post="/api/resources/{resource_id}/report"),
    GetResourceResp getResource(1: GetResourceReq req)(api.get="/api/resources/{resource_id}"),
    
    // 资源评分相关API
    SubmitResourceRatingResp submitResourceRating(1: SubmitResourceRatingReq req)(api.post="/api/resource_ratings/{rating_id}"),
    DeleteResourceRatingResp deleteResourceRating(1: DeleteResourceRatingReq req)(api.delete="/api/resource_ratings/{rating_id}"),
    
    // 资源评论相关API
    SubmitResourceCommentResp submitResourceComment(1: SubmitResourceCommentReq req)(api.post="/api/resource_comments/{comment_id}"),
    DeleteResourceCommentResp deleteResourceComment(1: DeleteResourceCommentReq req)(api.delete="/api/resources_comments/{comment_id}"),
    GetResourceCommentsResp getResourceComments(1: GetResourceCommentsReq req)(api.get="/api/resource/{resource_id}/comment"),
}