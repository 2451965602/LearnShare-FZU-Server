namespace go module

struct BaseResp{
    required i32 code,
    required string message,
}

struct User{
    required i64 userId,
    required string username,
    optional string password,
    required string email,
    required i64 college_id,
    required i64 major_id,
    required string avatar_url,
    required i64 reputation_score,
    required i64 roleId,
    required string status,
    required i64 created_at,
    required i64 updated_at,
}

struct Course{
    required i64 courseId,
    required string courseName,
    required i64 teacherId,
    required double credit,
    required i64 majorId,
    required string grade,
    required string description,
    required i64 createdAt,
    required i64 updatedAt,
}

struct CourseRating {
    required i64 ratingId,
    required i64 userId,
    required i64 courseId,
    required i32 usefulness,          
    required bool isVisible,
    required i64 createdAt,
}

struct CourseComment {
    required i64 commentId,
    required i64 userId,
    required i64 courseId,
    required string content,
    optional i64 parentId,            // ✅ 允许 NULL → optional
    required i64 likes,
    required bool isVisible,
    required string status,
    required i64 createdAt,
}

// 如果有嵌套结构，也同步更新
struct CourseCommentWithUser {
    required i64 commentId,
    required User user,
    required i64 courseId,
    required string content,
    optional i64 parentId,            // ✅ 保持一致
    required i64 likes,
    required bool isVisible,
    required string status,
    required i64 createdAt,
}
struct ResourceTag {
    required i64 tagId,
    required string tagName,
}

struct Resource {
    required i64 resourceId,
    required string title,              // 资源标题
    optional string description,        // 资源描述
    required string fileType,           // 文件类型 (.pdf, .docx, .pptx, .zip)
    required i64 fileSize,             // 文件大小 (bytes)
    required i64 uploaderId,           // 上传者ID
    required i64 courseId,             // 关联课程ID
    required i64 downloadCount,        // 下载次数
    required double averageRating,     // 平均评分
    required i64 ratingCount,          // 评分数量
    required i32 status,               // 资源状态 (0:待审核, 1:已发布, 2:已拒绝)
    required i64 createdAt,            // 创建时间
    optional list<ResourceTag> tags,   // 资源标签
}

enum ResourceCommentStatus {
    NORMAL = 0,
    DELETED_BY_USER = 1,
    DELETED_BY_ADMIN = 2,
}

struct ResourceRating{
    required i64 ratingId,
    required i64 userId,
    required i64 resourceId,
    required double recommendation,
    required bool isVisible,
    required i64 createdAt,
}


struct ResourceComment{
    required i64 commentId,
    required i64 userId,
    required i64 resourceId,
    required string content,
    required i64 parentId,
    required i64 likes,
    required bool isVisible,
    required ResourceCommentStatus status,
    required i64 createdAt,
}

struct ResourceCommentWithUser{
    required i64 commentId,
    required User user,
    required i64 resourceId,
    required string content,
    required i64 parentId,
    required i64 likes,
    required bool isVisible,
    required ResourceCommentStatus status,
    required i64 createdAt,
}

struct College{
    required i64 collegeId,
    required string collegeName,
    required string school
}

struct Major{
    required i64 majorId,
    required string majorName,
    required i64 collegeId,
}

struct Teacher{
    required i64 teacherId,
    required string teacherName,
    required i64 collegeId
    required string introduction,
    required string email,
    required string avatar_url,
    required i64 created_at,
    required i64 updated_at,
}

struct Permission{
    required i64 permissionId,
    required string permissionName,
    required string description,
}

struct Role{
    required i64 roleId,
    required string roleName,
    required string description,
    required list<Permission> permissions,
}

struct RolePermission{
    required i64 roleId,
    required i64 permissionId,
}

struct Favorite{
    required i64 favoriteId,
    required i64 userId,
    required i64 targetId,
    required string targetType,
    required i64 createdAt,
}

struct review{
    required i64 reviewId,
    required i64 reviewerId,
    required i64 reporterId,
    required i64 targetId,
    required string targetType,
    required string reason,
    required string status,
    required i64 priority,
    required i64 createdAt,
}
