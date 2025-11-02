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
    required i64 collegeId,
    required i64 majorId,
    required string avatarUrl,
    required i64 reputationScore,
    required i64 roleId,
    required string status,
    required i64 createdAt,
    required i64 updatedAt,
}

struct Course{
    required i64 courseId,
    required string courseName,
    required i64 teacherId,
    required i64 credit,
    required i64 majorId,
    required string grade,
    required string description,
    required i64 createdAt,
    required i64 updatedAt,
}

struct CourseRating{
    required i64 ratingId,
    required i64 userId,
    required i64 courseId,
    required i64 recommendation,
    required string difficulty,
    required i64 workload,
    required i64 usefulness,
    required bool isVisible,
    required i64 createdAt,
}


struct CourseComment{
    required i64 commentId,
    required i64 userId,
    required i64 courseId,
    required string content,
    required i64 parentId,
    required i64 likes,
    required bool isVisible,
    required i64 status,
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
    required string filePath,           // 文件路径
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
    required i64 recommendation,
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

