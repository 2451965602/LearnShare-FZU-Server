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
    required string action, // "approve" or "reject"
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
    required i64 course_id(api.path="course_id"),
    required string action, // "approve" or "reject"
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

struct AuditCommentReq{
    required i64 comment_id(api.path="comment_id"),
    required string action, // "approve" or "reject"
}
struct AuditCommentResp{
    required model.BaseResp base_resp,
}

service AdminAuditService {
    GetResourceAuditListResp GetResourceAuditList(1:GetResourceAuditListReq req)(api.get="/api/admin/audit/resources"),
    AuditResourceResp AuditResource(1:AuditResourceReq req)(api.put="/api/admin/audit/resources/:review_id"),
    GetCourseAuditListResp GetCourseAuditList(1:GetCourseAuditListReq req)(api.get="/api/admin/audit/courses"),
    AuditCourseResp AuditCourse(1:AuditCourseReq req)(api.post="/api/admin/audit/courses/:review_id"),
    GetCommentAuditListResp GetCommentAuditList(1:GetCommentAuditListReq req)(api.get="/api/admin/audit/comments"),
    AuditCommentResp AuditComment(1:AuditCommentReq req)(api.post="/api/admin/audit/comments/:review_id"),
}







