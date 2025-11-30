namespace go audit
include "model.thrift"


struct GetResourceAuditListReq{
    required i32 page_num,
    required i32 page_size,
}
struct GetResourceAuditListResp{
    required model.BaseResp base_resp,
    required list<model.review> resource_review_list,
}

struct AuditResourceReq{
    required i64 review_id(api.path="review_id"),
    required string action,
}
struct AuditResourceResp{
    required model.BaseResp base_resp,
}


struct GetCourseAuditListReq{
    required i32 page_num,
    required i32 page_size,
}
struct GetCourseAuditListResp{
    required model.BaseResp base_resp,
    required list<model.Course> course_audit_list,
}

struct AuditCourseReq{
    required i64 review_id(api.path="review_id"),
    required string action,
}
struct AuditCourseResp{
    required model.BaseResp base_resp,
}

struct GetCommentAuditListReq{
    required i32 page_num,
    required i32 page_size,
}
struct GetCommentAuditListResp{
    required model.BaseResp base_resp,
    required list<model.CourseComment> comment_audit_list,
}

// 获取待审核课程评论列表
struct GetCourseCommentAuditListReq{
    required i32 page_num,
    required i32 page_size,
}
struct GetCourseCommentAuditListResp{
    required model.BaseResp base_resp,
    required list<model.CourseComment> comment_audit_list,
}

// 获取待审核资源评论列表
struct GetResourceCommentAuditListReq{
    required i32 page_num,
    required i32 page_size,
}
struct GetResourceCommentAuditListResp{
    required model.BaseResp base_resp,
    required list<model.review> comment_audit_list,
}

struct AuditCourseCommentReq{
    required i64 review_id(api.path="review_id"),
    required string action,
}
struct AuditCourseCommentResp{
    required model.BaseResp base_resp,
}

struct AuditResourceCommentReq{
    required i64 review_id(api.path="review_id"),
    required string action,
}
struct AuditResourceCommentResp{
    required model.BaseResp base_resp,
}

service AdminAuditService {
    GetResourceAuditListResp GetResourceAuditList(1:GetResourceAuditListReq req)(api.get="/api/admin/audit/resources"),
    AuditResourceResp AuditResource(1:AuditResourceReq req)(api.post="/api/admin/audit/resources/:review_id"),
    GetCourseAuditListResp GetCourseAuditList(1:GetCourseAuditListReq req)(api.get="/api/admin/audit/courses"),
    AuditCourseResp AuditCourse(1:AuditCourseReq req)(api.post="/api/admin/audit/courses/:review_id"),
    GetCommentAuditListResp GetCommentAuditList(1:GetCommentAuditListReq req)(api.get="/api/admin/audit/comments"),
    GetCourseCommentAuditListResp GetCourseCommentAuditList(1:GetCourseCommentAuditListReq req)(api.get="/api/admin/audit/course_comments"),
    GetResourceCommentAuditListResp GetResourceCommentAuditList(1:GetResourceCommentAuditListReq req)(api.get="/api/admin/audit/resource_comments"),
    AuditCourseCommentResp AuditCourseComment(1:AuditCourseCommentReq req)(api.post="/api/admin/audit/course_comments/:review_id"),
    AuditResourceCommentResp AuditResourceComment(1:AuditResourceCommentReq req)(api.post="/api/admin/audit/resource_comments/:review_id"),
}



