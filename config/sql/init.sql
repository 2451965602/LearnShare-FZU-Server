-- 智汇选课系统数据库表创建语句 (最终优化版V5 - 融合内嵌外键、索引精简与事务控制)
-- 生成时间：2025-11-07
-- ------------------------------------------------------------
-- 禁用外键检查，并设置字符集
SET
    @OLD_FOREIGN_KEY_CHECKS = @ @FOREIGN_KEY_CHECKS;

SET
    FOREIGN_KEY_CHECKS = 0;

SET
    NAMES utf8mb4;

-- 开启事务，确保DDL和DML的原子性
START TRANSACTION;

-- 通用选项：使用 InnoDB、utf8mb4、ROW_FORMAT=DYNAMIC
-- ----------------------------
-- 权限表 (permissions)
-- ----------------------------
DROP TABLE IF EXISTS `permissions`;

CREATE TABLE `permissions` (
    `permission_id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '权限ID',
    `permission_name` VARCHAR(50) NOT NULL COMMENT '权限名称',
    `description` VARCHAR(500) DEFAULT NULL COMMENT '权限描述',
    PRIMARY KEY (`permission_id`),
    UNIQUE KEY `uk_permission_name` (`permission_name`)
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '权限表';

-- ----------------------------
-- 角色表 (roles) - 使用 SMALLINT 节省空间
-- ----------------------------
DROP TABLE IF EXISTS `roles`;

CREATE TABLE `roles` (
    `role_id` SMALLINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '角色ID',
    `role_name` VARCHAR(50) NOT NULL COMMENT '角色名称',
    `description` VARCHAR(500) DEFAULT NULL COMMENT '角色描述',
    PRIMARY KEY (`role_id`),
    UNIQUE KEY `uk_role_name` (`role_name`)
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '角色表';

-- ----------------------------
-- 角色-权限关联表 (role_permissions) - 内嵌外键
-- ----------------------------
DROP TABLE IF EXISTS `role_permissions`;

CREATE TABLE `role_permissions` (
    `role_id` SMALLINT UNSIGNED NOT NULL COMMENT '角色ID',
    `permission_id` INT UNSIGNED NOT NULL COMMENT '权限ID',
    PRIMARY KEY (`role_id`, `permission_id`),
    KEY `idx_rp_permission` (`permission_id`),
    CONSTRAINT `fk_rp_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`role_id`) ON DELETE CASCADE,
    CONSTRAINT `fk_rp_permission` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`permission_id`) ON DELETE CASCADE
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '角色-权限关联表';

-- ----------------------------
-- 学院表 (colleges)
-- ----------------------------
DROP TABLE IF EXISTS `colleges`;

CREATE TABLE `colleges` (
    `college_id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '学院ID',
    `college_name` VARCHAR(100) NOT NULL COMMENT '学院名称',
    `school` VARCHAR(100) DEFAULT NULL COMMENT '所属学校',
    PRIMARY KEY (`college_id`),
    UNIQUE KEY `uk_college_name` (`college_name`)
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '学院表';

-- ----------------------------
-- 专业表 (majors) - 移除冗余索引 idx_major_name
-- ----------------------------
DROP TABLE IF EXISTS `majors`;

CREATE TABLE `majors` (
    `major_id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '专业ID',
    `major_name` VARCHAR(100) NOT NULL COMMENT '专业名称',
    `college_id` INT UNSIGNED NOT NULL COMMENT '学院ID',
    PRIMARY KEY (`major_id`),
    UNIQUE KEY `uk_major_college` (`major_name`, `college_id`),
    -- 联合唯一索引
    CONSTRAINT `fk_major_college` FOREIGN KEY (`college_id`) REFERENCES `colleges` (`college_id`) ON DELETE RESTRICT
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '专业表';

-- ----------------------------
-- 教师表 (teachers) - 教师简介使用 TEXT, 外键 ON DELETE SET NULL
-- ----------------------------
DROP TABLE IF EXISTS `teachers`;

CREATE TABLE `teachers` (
    `teacher_id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '教师ID',
    `name` VARCHAR(100) NOT NULL COMMENT '教师姓名',
    `college_id` INT UNSIGNED DEFAULT NULL COMMENT '学院ID (设为可空以支持 ON DELETE SET NULL)',
    `introduction` TEXT COMMENT '教师简介',
    `email` VARCHAR(100) DEFAULT NULL COMMENT '教师邮箱',
    `avatar_url` VARCHAR(255) DEFAULT NULL COMMENT '头像URL',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`teacher_id`),
    KEY `idx_teacher_name` (`name`),
    CONSTRAINT `fk_teacher_college` FOREIGN KEY (`college_id`) REFERENCES `colleges` (`college_id`) ON DELETE
    SET
        NULL
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '教师表';

-- ----------------------------
-- 用户表 (users) - 内嵌外键
-- ----------------------------
DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
    `user_id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户唯一标识',
    `username` VARCHAR(32) NOT NULL COMMENT '用户名',
    `password_hash` VARCHAR(255) NOT NULL COMMENT '加密密码',
    `email` VARCHAR(100) NOT NULL COMMENT '邮箱',
    `college_id` INT UNSIGNED DEFAULT NULL COMMENT '所属学院ID',
    `major_id` INT UNSIGNED DEFAULT NULL COMMENT '所属专业ID',
    `avatar_url` VARCHAR(255) DEFAULT NULL COMMENT '头像URL',
    `reputation_score` INT NOT NULL DEFAULT 80 COMMENT '信誉分',
    `role_id` SMALLINT UNSIGNED NOT NULL COMMENT '角色ID',
    -- SMALLINT
    `status` ENUM('active', 'inactive', 'locked', 'banned') NOT NULL DEFAULT 'inactive' COMMENT '账户状态 (新增banned)',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`user_id`),
    UNIQUE KEY `uk_username` (`username`),
    UNIQUE KEY `uk_email` (`email`),
    KEY `idx_user_college_major` (`college_id`, `major_id`),
    CONSTRAINT `fk_user_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`role_id`) ON DELETE RESTRICT,
    CONSTRAINT `fk_user_college` FOREIGN KEY (`college_id`) REFERENCES `colleges` (`college_id`) ON DELETE
    SET
        NULL,
        CONSTRAINT `fk_user_major` FOREIGN KEY (`major_id`) REFERENCES `majors` (`major_id`) ON DELETE
    SET
        NULL
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '用户表';

-- ----------------------------
-- 课程表 (courses) - 调整学分为 DECIMAL(2,1) 更精确
-- ----------------------------
DROP TABLE IF EXISTS `courses`;

CREATE TABLE `courses` (
    `course_id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '课程ID',
    `course_name` VARCHAR(200) NOT NULL COMMENT '课程名称',
    `teacher_id` INT UNSIGNED DEFAULT NULL COMMENT '教师ID (设为可空以支持 ON DELETE SET NULL)',
    `credit` DECIMAL(2, 1) NOT NULL COMMENT '学分',
    `major_id` INT UNSIGNED NOT NULL COMMENT '专业ID',
    `grade` VARCHAR(20) NOT NULL COMMENT '年级',
    `description` TEXT COMMENT '课程简介',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`course_id`),
    KEY `idx_course_name` (`course_name`),
    KEY `idx_course_teacher_major` (`teacher_id`, `major_id`),
    CONSTRAINT `fk_course_teacher` FOREIGN KEY (`teacher_id`) REFERENCES `teachers` (`teacher_id`) ON DELETE
    SET
        NULL,
        CONSTRAINT `fk_course_major` FOREIGN KEY (`major_id`) REFERENCES `majors` (`major_id`) ON DELETE RESTRICT
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '课程表';

-- ----------------------------
-- 标签表 (tags)
-- ----------------------------
DROP TABLE IF EXISTS `tags`;

CREATE TABLE `tags` (
    `tag_id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '标签ID',
    `tag_name` VARCHAR(50) NOT NULL COMMENT '标签名称',
    PRIMARY KEY (`tag_id`),
    UNIQUE KEY `uk_tag_name` (`tag_name`)
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '标签表';

-- ----------------------------
-- 资源表 (resources) - 内嵌外键
-- ----------------------------
DROP TABLE IF EXISTS `resources`;

CREATE TABLE `resources` (
    `resource_id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '资源ID',
    `resource_name` VARCHAR(255) NOT NULL COMMENT '资源标题',
    `description` VARCHAR(500) DEFAULT NULL COMMENT '资源简介',
    `resource_url` VARCHAR(255) NOT NULL COMMENT '文件URL',
    `type` ENUM('pdf', 'docx', 'pptx', 'zip', 'other') NOT NULL COMMENT '文件类型 (新增other)',
    `size` INT UNSIGNED NOT NULL COMMENT '文件大小(字节)',
    `uploader_id` INT UNSIGNED NOT NULL COMMENT '上传者ID',
    `course_id` INT UNSIGNED DEFAULT NULL COMMENT '课程ID (允许NULL)',
    `download_count` INT UNSIGNED DEFAULT 0 COMMENT '下载次数',
    `average_rating` DECIMAL(2, 1) DEFAULT 0.0 COMMENT '平均评分',
    `rating_count` INT UNSIGNED DEFAULT 0 COMMENT '评分人数',
    `status` ENUM(
        'normal',
        'low_quality',
        'pending_review',
        'banned'
    ) DEFAULT 'pending_review' COMMENT '状态 (新增banned)',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`resource_id`),
    KEY `idx_resource_status` (`status`),
    KEY `idx_res_uploader_type` (`uploader_id`, `type`),
    KEY `idx_res_course_status` (`course_id`, `status`),
    CONSTRAINT `fk_resource_uploader` FOREIGN KEY (`uploader_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE,
    CONSTRAINT `fk_resource_course` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE
    SET
        NULL
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '资源表';

-- ----------------------------
-- 资源标签关联表 (resource_tags) - 内嵌外键
-- ----------------------------
DROP TABLE IF EXISTS `resource_tags`;

CREATE TABLE `resource_tags` (
    `resource_id` INT UNSIGNED NOT NULL COMMENT '资源ID',
    `tag_id` INT UNSIGNED NOT NULL COMMENT '标签ID',
    PRIMARY KEY (`resource_id`, `tag_id`),
    KEY `fk_rt_tag` (`tag_id`),
    CONSTRAINT `fk_rt_resource` FOREIGN KEY (`resource_id`) REFERENCES `resources` (`resource_id`) ON DELETE CASCADE,
    CONSTRAINT `fk_rt_tag` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`tag_id`) ON DELETE CASCADE
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '资源标签关联表';

-- ----------------------------
-- 课程评分表 (course_ratings) - 评分子项使用 TINYINT
-- ----------------------------
DROP TABLE IF EXISTS `course_ratings`;

CREATE TABLE `course_ratings` (
    `rating_id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '评分ID',
    `user_id` INT UNSIGNED NOT NULL COMMENT '用户ID',
    `course_id` INT UNSIGNED NOT NULL COMMENT '课程ID',
    `recommendation` DECIMAL(2, 1) NOT NULL COMMENT '综合推荐度(0.0-5.0)',
    `difficulty` TINYINT UNSIGNED NOT NULL COMMENT '课程难度(1-5)',
    `workload` TINYINT UNSIGNED NOT NULL COMMENT '作业压力(1-5)',
    `usefulness` TINYINT UNSIGNED NOT NULL COMMENT '知识实用性(1-5)',
    `is_visible` BOOLEAN DEFAULT TRUE COMMENT '是否显示',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`rating_id`),
    UNIQUE KEY `uk_user_course` (`user_id`, `course_id`),
    KEY `idx_cr_course_visible` (`course_id`, `is_visible`),
    CONSTRAINT `fk_cr_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE,
    CONSTRAINT `fk_cr_course` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '课程评分表';

-- ----------------------------
-- 课程评论表 (course_comments) - 内嵌外键
-- ----------------------------
DROP TABLE IF EXISTS `course_comments`;

CREATE TABLE `course_comments` (
    `comment_id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '评论ID',
    `user_id` INT UNSIGNED NOT NULL COMMENT '用户ID',
    `course_id` INT UNSIGNED NOT NULL COMMENT '课程ID',
    `content` TEXT NOT NULL COMMENT '评论内容',
    `parent_id` INT UNSIGNED DEFAULT NULL COMMENT '父评论ID',
    `likes` INT UNSIGNED DEFAULT 0 COMMENT '点赞数',
    `is_visible` BOOLEAN DEFAULT TRUE COMMENT '是否显示',
    `status` ENUM('normal', 'deleted_by_user', 'deleted_by_admin') DEFAULT 'normal' COMMENT '评论状态',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`comment_id`),
    KEY `fk_cc_parent` (`parent_id`),
    KEY `idx_cc_course_status` (`course_id`, `status`),
    CONSTRAINT `fk_cc_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE,
    CONSTRAINT `fk_cc_course` FOREIGN KEY (`course_id`) REFERENCES `courses` (`course_id`) ON DELETE CASCADE,
    CONSTRAINT `fk_cc_parent` FOREIGN KEY (`parent_id`) REFERENCES `course_comments` (`comment_id`) ON DELETE
    SET
        NULL
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '课程评论表';

-- ----------------------------
-- 课程评论反应表 (course_comment_reactions)
-- ----------------------------
DROP TABLE IF EXISTS `course_comment_reactions`;

CREATE TABLE `course_comment_reactions` (
    `reaction_id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '反应ID',
    `user_id` INT UNSIGNED NOT NULL COMMENT '用户ID',
    `comment_id` INT UNSIGNED NOT NULL COMMENT '课程评论ID',
    `reaction` ENUM('like', 'dislike') NOT NULL COMMENT '反应类型',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`reaction_id`),
    UNIQUE KEY `uk_ccr_user_comment` (`user_id`, `comment_id`) COMMENT '确保用户对同一评论仅有一个反应',
    KEY `idx_ccr_comment_reaction` (`comment_id`, `reaction`) COMMENT '加速按评论和反应类型查询',
    CONSTRAINT `fk_ccr_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE,
    CONSTRAINT `fk_ccr_comment` FOREIGN KEY (`comment_id`) REFERENCES `course_comments` (`comment_id`) ON DELETE CASCADE
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '课程评论反应表（点赞/点踩）';

-- ----------------------------
-- 资源评分表 (resource_ratings)
-- ----------------------------
DROP TABLE IF EXISTS `resource_ratings`;

CREATE TABLE `resource_ratings` (
    `rating_id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '评分ID',
    `user_id` INT UNSIGNED NOT NULL COMMENT '用户ID',
    `resource_id` INT UNSIGNED NOT NULL COMMENT '资源ID',
    `recommendation` DECIMAL(2, 1) NOT NULL COMMENT '综合推荐度(0.0-5.0)',
    `is_visible` BOOLEAN DEFAULT TRUE COMMENT '是否显示',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`rating_id`),
    UNIQUE KEY `uk_user_resource` (`user_id`, `resource_id`),
    KEY `idx_rr_resource_visible` (`resource_id`, `is_visible`),
    CONSTRAINT `fk_rr_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE,
    CONSTRAINT `fk_rr_resource` FOREIGN KEY (`resource_id`) REFERENCES `resources` (`resource_id`) ON DELETE CASCADE
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '资源评分表';

-- ----------------------------
-- 资源评论表 (resource_comments)
-- ----------------------------
DROP TABLE IF EXISTS `resource_comments`;

CREATE TABLE `resource_comments` (
    `comment_id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '评论ID',
    `user_id` INT UNSIGNED NOT NULL COMMENT '用户ID',
    `resource_id` INT UNSIGNED NOT NULL COMMENT '资源ID',
    `content` TEXT NOT NULL COMMENT '评论内容',
    `parent_id` INT UNSIGNED DEFAULT NULL COMMENT '父评论ID',
    `likes` INT UNSIGNED DEFAULT 0 COMMENT '点赞数',
    `is_visible` BOOLEAN DEFAULT TRUE COMMENT '是否显示',
    `status` ENUM('normal', 'deleted_by_user', 'deleted_by_admin') DEFAULT 'normal' COMMENT '评论状态',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`comment_id`),
    KEY `fk_rc_parent` (`parent_id`),
    KEY `idx_rc_resource_status` (`resource_id`, `status`),
    CONSTRAINT `fk_rc_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE,
    CONSTRAINT `fk_rc_resource` FOREIGN KEY (`resource_id`) REFERENCES `resources` (`resource_id`) ON DELETE CASCADE,
    CONSTRAINT `fk_rc_parent` FOREIGN KEY (`parent_id`) REFERENCES `resource_comments` (`comment_id`) ON DELETE
    SET
        NULL
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '资源评论表';

-- ----------------------------
-- 资源评论反应表 (resource_comment_reactions)
-- ----------------------------
DROP TABLE IF EXISTS `resource_comment_reactions`;

CREATE TABLE `resource_comment_reactions` (
    `reaction_id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` INT UNSIGNED NOT NULL,
    `comment_id` INT UNSIGNED NOT NULL,
    `reaction` ENUM('like', 'dislike') NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`reaction_id`),
    UNIQUE KEY `uk_rcr_user_comment` (`user_id`, `comment_id`),
    KEY `idx_rcr_comment_reaction` (`comment_id`, `reaction`),
    CONSTRAINT `fk_rcr_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE,
    CONSTRAINT `fk_rcr_comment` FOREIGN KEY (`comment_id`) REFERENCES `resource_comments` (`comment_id`) ON DELETE CASCADE
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '资源评论反应表';

-- ----------------------------
-- 信誉分记录表 (reputation_records) - 调整 change_score 为 SMALLINT
-- ----------------------------
DROP TABLE IF EXISTS `reputation_records`;

CREATE TABLE `reputation_records` (
    `record_id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '记录ID',
    `user_id` INT UNSIGNED NOT NULL COMMENT '用户ID',
    `change_score` SMALLINT NOT NULL COMMENT '分数变化',
    `reason` VARCHAR(500) NOT NULL COMMENT '变化原因',
    `related_id` INT UNSIGNED DEFAULT NULL COMMENT '关联对象ID',
    `related_type` ENUM('resource', 'comment', 'rating') DEFAULT NULL COMMENT '关联对象类型',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录时间',
    PRIMARY KEY (`record_id`),
    KEY `idx_reputation_time` (`created_at`),
    CONSTRAINT `fk_reprec_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '信誉分记录表';

-- ----------------------------
-- 收藏表 (favorites)
-- ----------------------------
DROP TABLE IF EXISTS `favorites`;

CREATE TABLE `favorites` (
    `favorite_id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '收藏ID',
    `user_id` INT UNSIGNED NOT NULL COMMENT '用户ID',
    `target_id` INT UNSIGNED NOT NULL COMMENT '对象ID',
    `target_type` ENUM('course', 'resource', 'teacher') NOT NULL COMMENT '对象类型',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '收藏时间',
    PRIMARY KEY (`favorite_id`),
    UNIQUE KEY `uk_user_target` (`user_id`, `target_id`, `target_type`),
    KEY `idx_favorite_type` (`user_id`, `target_type`),
    CONSTRAINT `fk_favorite_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '收藏表';

-- ----------------------------
-- 审核表 (reviews) - 优化状态和优先级索引，新增举报人外键
-- ----------------------------
DROP TABLE IF EXISTS `reviews`;

CREATE TABLE `reviews` (
    `review_id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '审核ID',
    `target_id` INT UNSIGNED NOT NULL COMMENT '对象ID',
    `reporter_id` INT UNSIGNED NOT NULL COMMENT '举报人ID',
    `target_type` ENUM(
        'resource',
        'course_rating',
        'resource_rating',
        'comment'
    ) NOT NULL COMMENT '对象类型',
    `reason` VARCHAR(500) NOT NULL COMMENT '审核原因',
    `status` ENUM('pending', 'approved', 'rejected') DEFAULT 'pending' COMMENT '状态',
    `priority` TINYINT UNSIGNED DEFAULT 3 COMMENT '优先级(1-5)',
    `reviewer_id` INT UNSIGNED DEFAULT NULL COMMENT '审核员ID',
    `reviewed_at` TIMESTAMP NULL DEFAULT NULL COMMENT '审核时间',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`review_id`),
    KEY `idx_review_status_priority` (`status`, `priority`),
    -- 优化索引
    KEY `idx_review_target` (`target_type`, `target_id`),
    CONSTRAINT `fk_review_reporter` FOREIGN KEY (`reporter_id`) REFERENCES `users` (`user_id`) ON DELETE RESTRICT,
    CONSTRAINT `fk_review_reviewer` FOREIGN KEY (`reviewer_id`) REFERENCES `users` (`user_id`) ON DELETE
    SET
        NULL
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '审核表';

-- ----------------------------
-- 物品表 (items)
-- ----------------------------
DROP TABLE IF EXISTS `items`;

CREATE TABLE `items` (
    `item_id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '物品ID',
    `name` VARCHAR(100) NOT NULL COMMENT '物品名称',
    `type` ENUM(
        'avatar_frame',
        'background',
        'medal',
        'nickname_color'
    ) NOT NULL COMMENT '物品类型',
    `price` INT UNSIGNED NOT NULL COMMENT '所需积分',
    `description` VARCHAR(500) DEFAULT NULL COMMENT '物品描述',
    `image_url` VARCHAR(255) DEFAULT NULL COMMENT '预览图URL',
    PRIMARY KEY (`item_id`),
    UNIQUE KEY `uk_item_name` (`name`)
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '物品表';

-- ----------------------------
-- 用户物品表 (user_items) - 新增查询启用状态的索引
-- ----------------------------
DROP TABLE IF EXISTS `user_items`;

CREATE TABLE `user_items` (
    `user_item_id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '记录ID',
    `user_id` INT UNSIGNED NOT NULL COMMENT '用户ID',
    `item_id` INT UNSIGNED NOT NULL COMMENT '物品ID',
    `is_used` BOOLEAN DEFAULT FALSE COMMENT '是否启用 (使用BOOLEAN)',
    `obtained_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '获取时间',
    PRIMARY KEY (`user_item_id`),
    UNIQUE KEY `uk_user_item` (`user_id`, `item_id`),
    KEY `idx_ui_user_used` (`user_id`, `is_used`),
    CONSTRAINT `fk_ui_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE,
    CONSTRAINT `fk_ui_item` FOREIGN KEY (`item_id`) REFERENCES `items` (`item_id`) ON DELETE CASCADE
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '用户物品表';

-- ----------------------------
-- 数据初始化: 角色 (Roles)
-- ----------------------------
INSERT INTO
    `roles` (`role_id`, `role_name`, `description`)
VALUES
    (1, '超级管理员', '拥有平台全部系统配置、账号管理与内容审核权限'),
    (2, '普通用户', '拥有基本功能'),
    (3, '审核员', '负责全站资源与评论的审核处理，以及用户管理辅助');

-- ----------------------------
-- 数据初始化: 权限 (Permissions)
-- ----------------------------
INSERT INTO
    `permissions` (
        `permission_id`,
        `permission_name`,
        `description`
    )
VALUES
    (1, 'user.profile.update', '修改个人基础资料与账户安全信息'),
    (2, 'resource.download', '下载学习资源文件'),
    (3, 'resource.upload', '上传并发布新的学习资源'),
    (4, 'resource.manage_all', '管理全站资源内容与状态'),
    (5, 'resource.comment.moderate', '审核、隐藏或删除资源评论'),
    (6, 'resource.rating.moderate', '审核或删除资源评分'),
    (7, 'course.comment.moderate', '审核或删除课程评论'),
    (8, 'course.rating.moderate', '审核或删除课程评分'),
    (9, 'review.handle', '处理举报、审核与处罚流程'),
    (10, 'role.manage', '维护角色、权限及其分配关系'),
    (11, 'user.reputation.adjust', '调整用户信誉分与记录变动原因'),
    (12, 'user.account.manage', '冻结、解封或禁用用户账户'),
    (13, 'content.tag.manage', '编辑或删除资源与课程标签'),
    (14, 'content.category.manage', '管理课程分类与专业信息'),
    (15, 'teacher.profile.manage', '维护教师信息与关联课程'),
    (16, 'system.config.manage', '修改系统配置与全局参数'),
    (17, 'audit.log.view', '查看系统操作日志与审核记录'),
    (18, 'data.report.view', '访问平台数据报表与统计信息');

-- ----------------------------
-- 数据初始化: 角色权限映射 (Role Permissions)
-- ----------------------------
-- 审核员 (Role ID: 3)
INSERT INTO
    `role_permissions` (`role_id`, `permission_id`)
VALUES
    (3, 4),
    (3, 5),
    (3, 6),
    (3, 7),
    (3, 8),
    (3, 9),
    (3, 11),
    (3, 12),
    (3, 13),
    (3, 14),
    (3, 15),
    (3, 17),
    (3, 18);

-- 超级管理员 (Role ID: 1)
INSERT INTO
    `role_permissions` (`role_id`, `permission_id`)
SELECT
    1,
    `permission_id`
FROM
    `permissions`;

-- 提交事务并恢复环境参数
COMMIT;

SET
    FOREIGN_KEY_CHECKS = @OLD_FOREIGN_KEY_CHECKS;