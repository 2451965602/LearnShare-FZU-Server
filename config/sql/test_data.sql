-- ==========================================================
-- 智汇选课系统 – 最终合并版测试数据（与 Optimized V5 兼容）
-- 生成日期：2025-11-07
-- 说明：
--   1. 确保在 V5 结构脚本执行成功后运行此脚本。
--   2. 密码统一为 admin123（散列值：$2a$10$g6AlcmLuZQLbxAh3DkRL7.xSnM0firkHssuC9g7LcfvfZpJH6n/8.）。
-- ==========================================================
-- 0. 关闭外键检查，允许自由插入
SET FOREIGN_KEY_CHECKS = 0;

-- ------------------------------------------------------------
-- 清空已有业务数据（可选：如果需要重新插入）
-- ------------------------------------------------------------
TRUNCATE TABLE `user_items`;
TRUNCATE TABLE `items`;
TRUNCATE TABLE `reviews`;
TRUNCATE TABLE `favorites`;
TRUNCATE TABLE `reputation_records`;
TRUNCATE TABLE `resource_comments`;
TRUNCATE TABLE `resource_comment_reactions`;
TRUNCATE TABLE `resource_ratings`;
TRUNCATE TABLE `course_comments`;
TRUNCATE TABLE `course_ratings`;
TRUNCATE TABLE `resource_tags`;
TRUNCATE TABLE `resources`;
TRUNCATE TABLE `tags`;
TRUNCATE TABLE `courses`;
TRUNCATE TABLE `users`;
TRUNCATE TABLE `teachers`;
TRUNCATE TABLE `majors`;
TRUNCATE TABLE `colleges`;

-- ----------------------------
-- 1. 学院 colleges
-- ----------------------------
INSERT INTO `colleges` (`college_name`, `school`) VALUES
                                                      ('计算机学院', '福州大学'),
                                                      ('商学院', '福州大学'),
                                                      ('理学院', '福州大学');

-- ----------------------------
-- 2. 专业 majors
-- ----------------------------
INSERT INTO `majors` (`major_name`, `college_id`) VALUES
                                                      ('计算机科学与技术', 1),
                                                      ('人工智能', 1),
                                                      ('金融学', 2),
                                                      ('数学与应用数学', 3);

-- ----------------------------
-- 3. 教师 teachers
-- ----------------------------
INSERT INTO `teachers` (`name`, `college_id`, `introduction`) VALUES
                                                                  ('张教授', 1, '研究方向：分布式系统、操作系统。擅长理论与实践结合教学。'),
                                                                  ('李教授', 1, '研究方向：机器学习、人工智能。指导过多个顶级算法竞赛团队。'),
                                                                  ('王教授', 2, '研究方向：金融科技、量化投资。拥有丰富的业界经验。'),
                                                                  ('赵教授', 3, '研究方向：运筹学、数理统计。教学严谨，基础扎实。');

-- ----------------------------
-- 4. 标签 tags
-- ----------------------------
INSERT INTO `tags` (`tag_name`) VALUES
                                    ('讲义'), ('习题'), ('历年试卷'), ('PPT'), ('实验报告'),
                                    ('复习资料'), ('课程设计'), ('项目代码'), ('参考书目'), ('视频教程'),
                                    ('考研资料'), ('面试题库'), ('开源项目'), ('学习笔记'), ('算法模板');

-- ----------------------------
-- 5. 用户 users
--    密码统一为 test123456 -> $2a$10$g6AlcmLuZQLbxAh3DkRL7.xSnM0firkHssuC9g7LcfvfZpJH6n/8.
-- ----------------------------
INSERT INTO `users` (`username`, `password_hash`, `email`, `college_id`, `major_id`, `role_id`, `status`, `reputation_score`) VALUES
                                                                                                                                  ('admin',    '$2a$10$g6AlcmLuZQLbxAh3DkRL7.xSnM0firkHssuC9g7LcfvfZpJH6n/8.', 'admin@example.com',    1, 1, 1, 'active', 100), -- ID: 1
                                                                                                                                  ('reviewer', '$2a$10$g6AlcmLuZQLbxAh3DkRL7.xSnM0firkHssuC9g7LcfvfZpJH6n/8.', 'reviewer@example.com', 1, 1, 3, 'active', 100), -- ID: 2
                                                                                                                                  ('alice',    '$2a$10$g6AlcmLuZQLbxAh3DkRL7.xSnM0firkHssuC9g7LcfvfZpJH6n/8.', 'alice@example.com',    1, 1, 2, 'active', 90),  -- ID: 3
                                                                                                                                  ('bob',      '$2a$10$g6AlcmLuZQLbxAh3DkRL7.xSnM0firkHssuC9g7LcfvfZpJH6n/8.', 'bob@example.com',      1, 2, 2, 'active', 85),  -- ID: 4
                                                                                                                                  ('carol',    '$2a$10$g6AlcmLuZQLbxAh3DkRL7.xSnM0firkHssuC9g7LcfvfZpJH6n/8.', 'carol@example.com',    2, 3, 2, 'active', 75),  -- ID: 5
                                                                                                                                  ('david',    '$2a$10$g6AlcmLuZQLbxAh3DkRL7.xSnM0firkHssuC9g7LcfvfZpJH6n/8.', 'david@example.com',    1, 1, 2, 'active', 95),  -- ID: 6
                                                                                                                                  ('emma',     '$2a$10$g6AlcmLuZQLbxAh3DkRL7.xSnM0firkHssuC9g7LcfvfZpJH6n/8.', 'emma@example.com',     1, 2, 2, 'active', 88),  -- ID: 7
                                                                                                                                  ('frank',    '$2a$10$g6AlcmLuZQLbxAh3DkRL7.xSnM0firkHssuC9g7LcfvfZpJH6n/8.', 'frank@example.com',    3, 4, 2, 'active', 82),  -- ID: 8
                                                                                                                                  ('grace',    '$2a$10$g6AlcmLuZQLbxAh3DkRL7.xSnM0firkHssuC9g7LcfvfZpJH6n/8.', 'grace@example.com',    2, 3, 2, 'inactive', 60), -- ID: 9
                                                                                                                                  ('henry',    '$2a$10$g6AlcmLuZQLbxAh3DkRL7.xSnM0firkHssuC9g7LcfvfZpJH6n/8.', 'henry@example.com',    1, 1, 2, 'locked', 45);  -- ID: 10

-- ----------------------------
-- 6. 课程 courses
-- ----------------------------
INSERT INTO `courses` (`course_name`, `teacher_id`, `credit`, `major_id`, `grade`, `description`) VALUES
                                                                                                      ('操作系统', 1, 3.0, 1, '大三', '深入理解操作系统原理，包括进程管理、内存管理、文件系统等核心概念'),   -- ID: 1
                                                                                                      ('深度学习', 2, 2.5, 2, '研一', '神经网络与框架实战，涵盖CNN、RNN、Transformer等主流架构'),          -- ID: 2
                                                                                                      ('数据结构与算法', 1, 4.0, 1, '大二', '程序设计基础，培养算法思维和编程能力'),                           -- ID: 3
                                                                                                      ('计算机网络', 1, 3.0, 1, '大三', 'TCP/IP协议栈详解，网络编程基础'),                                    -- ID: 4
                                                                                                      ('数据库系统', 2, 3.5, 1, '大三', '关系型数据库设计，SQL优化，NoSQL介绍'),                               -- ID: 5
                                                                                                      ('机器学习基础', 2, 3.0, 2, '大四', '监督学习、无监督学习理论基础与实践'),                                -- ID: 6
                                                                                                      ('金融工程', 3, 3.0, 3, '大四', '衍生品定价与风险管理，量化投资基础'),                                   -- ID: 7
                                                                                                      ('市场营销', 3, 2.5, 3, '大二', '消费者行为分析，品牌管理策略'),                                           -- ID: 8
                                                                                                      ('财务管理', 3, 3.0, 3, '大三', '企业财务分析，投资决策方法'),                                           -- ID: 9
                                                                                                      ('高等数学', 4, 4.0, 4, '大一', '微积分与线性代数基础'),                                                  -- ID: 10
                                                                                                      ('概率论与数理统计', 4, 3.0, 4, '大二', '随机变量理论，假设检验，回归分析'),                              -- ID: 11
                                                                                                      ('离散数学', 4, 3.0, 4, '大二', '逻辑推理，图论，组合数学基础');                                          -- ID: 12

-- ----------------------------
-- 7. 资源 resources
-- ----------------------------
INSERT INTO `resources` (`resource_name`, `description`, `resource_url`, `type`, `size`, `uploader_id`, `course_id`, `status`) VALUES
                                                                                                                                   ('操作系统-完整讲义', '张教授操作系统课程完整讲义，涵盖进程管理、内存管理、文件系统等章节', 'https://file.example.com/os-complete-notes.pdf', 'pdf', 5248000, 3, 1, 'normal'),  -- ID: 1
                                                                                                                                   ('OS-Lab实验代码', '操作系统实验课代码，包括进程调度算法实现', 'https://file.example.com/os-lab-code.zip', 'zip', 1024000, 4, 1, 'normal'),  -- ID: 2
                                                                                                                                   ('操作系统-历年期末试卷', '2020-2024年期末试卷及答案解析', 'https://file.example.com/os-exam-papers.pdf', 'pdf', 3072000, 5, 1, 'normal'),  -- ID: 3
                                                                                                                                   ('深度学习-PyTorch实战', 'PyTorch框架入门到精通，包含大量实例代码', 'https://file.example.com/dl-pytorch-tutorial.zip', 'zip', 8388608, 3, 2, 'normal'),  -- ID: 4
                                                                                                                                   ('CNN经典论文合集', 'LeNet、AlexNet、VGG、ResNet等经典论文及代码实现', 'https://file.example.com/cnn-papers.zip', 'zip', 15728640, 6, 2, 'normal'), -- ID: 5
                                                                                                                                   ('深度学习面试题库', '2024年最新深度学习面试题及答案', 'https://file.example.com/dl-interview-questions.pdf', 'pdf', 2048000, 7, 2, 'pending_review'), -- ID: 6
                                                                                                                                   ('算法竞赛模板', 'ACM算法竞赛常用模板，包括图论、动态规划等', 'https://file.example.com/algorithm-templates.zip', 'zip', 512000, 3, 3, 'normal'),  -- ID: 7
                                                                                                                                   ('LeetCode刷题笔记', '精选300道LeetCode题目详解与代码实现', 'https://file.example.com/leetcode-notes.pdf', 'pdf', 4096000, 4, 3, 'normal'),  -- ID: 8
                                                                                                                                   ('量化投资策略回测', '基于Python的量化投资策略及回测系统', 'https://file.example.com/quant-strategy.zip', 'zip', 6291456, 8, 7, 'normal'),  -- ID: 9
                                                                                                                                   ('期权定价模型实现', 'Black-Scholes、二叉树等期权定价模型代码', 'https://file.example.com/option-pricing.docx', 'docx', 1536000, 9, 7, 'low_quality'), -- ID: 10
                                                                                                                                   ('高等数学-考研复习资料', '2025考研数学一复习资料及历年真题', 'https://file.example.com/math-exam-review.zip', 'zip', 20480000, 10, 10, 'normal'), -- ID: 11
                                                                                                                                   ('概率统计重点公式总结', '概率论与数理统计重要公式及考点梳理', 'https://file.example.com/prob-formulas.pdf', 'pdf', 1024000, 3, 11, 'normal'); -- ID: 12

-- ----------------------------
-- 8. 资源标签 resource_tags
-- ----------------------------
INSERT INTO `resource_tags` (`resource_id`, `tag_id`) VALUES
                                                          (1, 1), (1, 6), (2, 5), (2, 8), (3, 3), (3, 6),
                                                          (4, 1), (4, 10), (5, 8), (5, 13), (6, 12), (6, 6),
                                                          (7, 15), (7, 8), (8, 1), (8, 6), (9, 8), (9, 7),
                                                          (10, 1), (11, 11), (11, 6), (12, 1), (12, 14);

-- ----------------------------
-- 9. 课程评分 course_ratings
-- ----------------------------
INSERT INTO `course_ratings` (`user_id`, `course_id`, `recommendation`, `difficulty`, `workload`, `usefulness`) VALUES
                                                                                                                    (3, 1, 4.5, 3, 2, 5), (4, 1, 3.5, 4, 3, 4), (5, 1, 4.8, 3, 2, 5), (6, 1, 4.2, 3, 3, 4), -- OS
                                                                                                                    (3, 2, 5.0, 5, 4, 5), (7, 2, 4.7, 4, 4, 5), (8, 2, 4.0, 5, 5, 4),                       -- 深度学习
                                                                                                                    (4, 3, 4.3, 3, 4, 5), (6, 3, 4.9, 2, 3, 5), (9, 3, 3.8, 4, 4, 4),                       -- 数据结构
                                                                                                                    (5, 4, 4.1, 3, 3, 4), (7, 4, 3.9, 4, 3, 4),                                           -- 计算机网络
                                                                                                                    (3, 5, 4.6, 3, 3, 5), (8, 5, 4.4, 3, 4, 4),                                           -- 数据库系统
                                                                                                                    (8, 7, 4.2, 4, 3, 4), (9, 7, 3.7, 5, 4, 3),                                           -- 金融工程
                                                                                                                    (4, 10, 3.9, 4, 4, 4), (5, 10, 4.5, 3, 3, 5), (10, 10, 3.2, 5, 5, 3);                   -- 高等数学

-- ----------------------------
-- 10. 课程评论 course_comments
-- ----------------------------
INSERT INTO `course_comments` (`user_id`, `course_id`, `content`, `parent_id`, `likes`, `status`) VALUES
                                                                                                      (3, 1, '张教授讲课非常清晰！理论与实践结合得很好，特别是进程同步的章节，收获很大。', NULL, 15, 'normal'), -- ID: 1
                                                                                                      (4, 1, '作业量略大，但收获很多，建议多安排一些实验课时间', 1, 8, 'normal'),                           -- ID: 2
                                                                                                      (5, 1, '老师的课件很详细，课后复习很方便', NULL, 6, 'normal'),                                         -- ID: 3
                                                                                                      (6, 1, '有没有同学愿意一起组队做OS的课程设计？', 3, 3, 'normal'),                                     -- ID: 4
                                                                                                      (3, 2, '李教授的深度学习课程太棒了！从基础理论到最新研究都覆盖到了', NULL, 20, 'normal'),              -- ID: 5
                                                                                                      (7, 2, '课程内容充实，但是确实比较难，需要一定的数学基础', 5, 12, 'normal'),                           -- ID: 6
                                                                                                      (8, 2, '实验环节很有趣，跑通了第一个CNN模型的时候特别有成就感', 5, 10, 'normal'),                       -- ID: 7
                                                                                                      (4, 2, '推荐提前看一些机器学习的基础课程', NULL, 5, 'normal'),                                         -- ID: 8
                                                                                                      (6, 3, '张教授的算法课讲得真好，复杂度分析部分豁然开朗', NULL, 18, 'normal'),                          -- ID: 9
                                                                                                      (9, 3, '课程设计很有挑战性，但是完成后感觉算法思维提升了很多', 9, 9, 'normal'),                          -- ID: 10 (修正 parent_id)
                                                                                                      (4, 3, '建议多做练习题，OJ平台上的题目很有帮助', NULL, 7, 'normal'),                                 -- ID: 11
                                                                                                      (3, 5, '李教授的数据库课程很实用，SQL优化部分特别有用', NULL, 14, 'normal'),                         -- ID: 12
                                                                                                      (8, 5, '项目作业很有意思，设计了一个小型图书管理系统', 12, 8, 'normal'),                              -- ID: 13 (修正 parent_id)
                                                                                                      (8, 7, '王教授的金融工程课程难度不小，但是内容很前沿', NULL, 11, 'normal'),                          -- ID: 14
                                                                                                      (9, 7, '量化交易策略回测实验很有趣，第一次感受到理论与实践的结合', 14, 6, 'normal'),                     -- ID: 15 (修正 parent_id)
                                                                                                      (4, 10, '赵教授的高数课很扎实，基础概念讲得很清楚', NULL, 16, 'normal'),                             -- ID: 16
                                                                                                      (5, 10, '课后习题有点多，但确实有助于理解概念', 16, 9, 'normal'),                                    -- ID: 17 (修正 parent_id)
                                                                                                      (10, 10, '希望老师能多讲一些考研的解题技巧', NULL, 4, 'normal');                                      -- ID: 18

-- ----------------------------
-- 11. 资源评分 resource_ratings
-- ----------------------------
INSERT INTO `resource_ratings` (`user_id`, `resource_id`, `recommendation`) VALUES
                                                                                (3, 1, 4.5), (4, 1, 4.8), (5, 1, 4.2), (6, 2, 4.6), (7, 2, 4.0), (4, 3, 4.9), -- OS
                                                                                (3, 4, 5.0), (8, 4, 4.7), (6, 5, 4.8), (7, 5, 4.4), (3, 6, 4.3),               -- DL
                                                                                (4, 7, 4.9), (6, 7, 4.7), (9, 7, 4.5), (3, 8, 4.6), (5, 8, 4.8),               -- 算法
                                                                                (8, 9, 4.2), (9, 9, 3.9),                                                      -- 金融
                                                                                (4, 11, 4.7), (5, 11, 4.5), (10, 11, 4.3), (3, 12, 4.4);                        -- 数学

-- ----------------------------
-- 12. 资源评论 resource_comments
-- ----------------------------
INSERT INTO `resource_comments` (`user_id`, `resource_id`, `content`, `parent_id`, `likes`) VALUES
                                                                                                (3, 1, '讲义非常详细，配合上课效果很好，感谢分享！', NULL, 12), -- ID: 1
                                                                                                (4, 1, '已下载，排版清晰，内容全面', 1, 6),                    -- ID: 2
                                                                                                (5, 1, '第三章的进程调度算法图解很直观', 1, 8),                -- ID: 3
                                                                                                (6, 2, '实验代码注释很详细，在Ubuntu上完美运行', NULL, 10),     -- ID: 4
                                                                                                (7, 2, 'FCFS和SJF算法的实现都很标准，谢谢学长', 4, 5),        -- ID: 5
                                                                                                (4, 3, '历年试卷很有参考价值，答案解析也很到位', NULL, 15),     -- ID: 6
                                                                                                (3, 4, 'PyTorch实战教程写得太好了，从零基础到上手项目全靠它', NULL, 20), -- ID: 7
                                                                                                (8, 4, '代码示例很丰富，每个模型都有完整的训练流程', 7, 12),    -- ID: 8
                                                                                                (6, 5, '经典论文合集很有价值，省去了自己到处找的时间', NULL, 18), -- ID: 9
                                                                                                (7, 5, 'ResNet的实现代码很有参考价值，直接复用了', 9, 9),      -- ID: 10
                                                                                                (3, 6, '面试题覆盖面很广，答案也很详细', NULL, 14),             -- ID: 11
                                                                                                (4, 7, '算法模板很实用，比赛时直接套用节省了很多时间', NULL, 16), -- ID: 12
                                                                                                (6, 7, '动态规划部分的模板特别好用，感谢分享', 12, 8),          -- ID: 13 (修正 parent_id)
                                                                                                (9, 7, '图论算法实现很规范，没有bug', 12, 7),                  -- ID: 14 (修正 parent_id)
                                                                                                (3, 8, 'LeetCode题目解析很详细，思路清晰', NULL, 13),           -- ID: 15
                                                                                                (5, 8, '300道题目精选得很好，都是经典题型', 15, 10),           -- ID: 16 (修正 parent_id)
                                                                                                (8, 9, '量化策略回测系统很专业，学到了很多', NULL, 11),         -- ID: 17
                                                                                                (9, 9, '代码结构清晰，容易理解和修改', 17, 6),                 -- ID: 18 (修正 parent_id)
                                                                                                (4, 11, '考研资料整理得很系统，知识点覆盖全面', NULL, 19),      -- ID: 19
                                                                                                (5, 11, '历年真题解析很详细，错误率应该不高', 19, 9),          -- ID: 20 (修正 parent_id)
                                                                                                (10, 11, '比市面上的一些考研资料要好很多', 19, 7),             -- ID: 21 (修正 parent_id)
                                                                                                (3, 12, '公式总结得很到位，考试前复习很方便', NULL, 8);         -- ID: 22

-- ----------------------------
-- 13. 信誉分记录 reputation_records
-- ----------------------------
INSERT INTO `reputation_records` (`user_id`, `change_score`, `reason`, `related_id`, `related_type`) VALUES
                                                                                                         (3, 15, '上传优质资源被下载50次，用户评价很高', 1, 'resource'),
                                                                                                         (3, 12, '上传的深度学习PyTorch教程获得高评分', 4, 'resource'),
                                                                                                         (3, 8, '多条课程评论获得点赞', 1, 'comment'),
                                                                                                         (4, 10, '上传的算法竞赛模板被广泛使用', 7, 'resource'),
                                                                                                         (4, 5, '课程评分被其他用户认为很有帮助', 1, 'rating'),
                                                                                                         (5, 8, '首次上传优质资源', 3, 'resource'),
                                                                                                         (6, 6, '资源评论获得点赞', 2, 'resource'),
                                                                                                         (7, 10, '上传的CNN经典论文合集受到好评', 5, 'resource'),
                                                                                                         (8, 7, '上传的量化策略回测系统代码被认可', 9, 'resource'),
                                                                                                         (4, 5, '上传的考研资料帮助了很多同学', 11, 'resource'),
                                                                                                         (9, -8, '上传的资源存在质量问题，被标记为low_quality', 10, 'resource'),
                                                                                                         (4, -3, '评论被其他用户举报并被审核员删除', 2, 'comment'),
                                                                                                         (5, -5, '评分被认为不够客观', 7, 'rating'),
                                                                                                         (3, 10, '完成账户设置并验证邮箱', 0, NULL),
                                                                                                         (4, 10, '完成账户设置并验证邮箱', 0, NULL),
                                                                                                         (5, 10, '完成账户设置并验证邮箱', 0, NULL),
                                                                                                         (6, 10, '完成账户设置并验证邮箱', 0, NULL),
                                                                                                         (7, 10, '完成账户设置并验证邮箱', 0, NULL),
                                                                                                         (8, 10, '完成账户设置并验证邮箱', 0, NULL),
                                                                                                         (9, 10, '完成账户设置并验证邮箱', 0, NULL),
                                                                                                         (10, 10, '完成账户设置并验证邮箱', 0, NULL);

-- ----------------------------
-- 14. 收藏 favorites
-- ----------------------------
INSERT INTO `favorites` (`user_id`, `target_id`, `target_type`) VALUES
                                                                    (3, 1, 'course'), (3, 2, 'course'), (3, 3, 'course'), (4, 1, 'course'), (4, 5, 'course'),
                                                                    (5, 2, 'course'), (5, 7, 'course'), (6, 3, 'course'), (6, 4, 'course'), (7, 2, 'course'),
                                                                    (8, 5, 'course'), (8, 7, 'course'),
                                                                    (3, 1, 'resource'), (3, 4, 'resource'), (3, 7, 'resource'), (4, 2, 'resource'), (4, 5, 'resource'),
                                                                    (4, 11, 'resource'), (5, 3, 'resource'), (5, 6, 'resource'), (6, 7, 'resource'), (6, 8, 'resource'),
                                                                    (7, 4, 'resource'), (7, 5, 'resource'),
                                                                    (3, 1, 'teacher'), (3, 2, 'teacher'), (4, 1, 'teacher'), (5, 3, 'teacher'), (6, 4, 'teacher');

-- ----------------------------
-- 15. 审核 reviews
-- ----------------------------
INSERT INTO `reviews` (`target_id`, `reporter_id`, `target_type`, `reason`, `priority`, `reviewer_id`, `status`) VALUES
                                                                                                                     (6, 3, 'resource', '深度学习面试题库需确认版权和内容准确性', 3, 2, 'pending'),
                                                                                                                     (10, 4, 'resource', '期权定价模型Excel文件可能有格式问题', 4, 2, 'pending'),
                                                                                                                     (2, 5, 'comment', '用户反馈课程讨论组信息过时，已更新', 2, 2, 'approved'),
                                                                                                                     (3, 6, 'resource', '历年试卷资源确认无版权问题', 3, 2, 'approved'),
                                                                                                                     (8, 7, 'comment', '课程评论中包含无关链接', 4, 2, 'rejected'),
                                                                                                                     (9, 8, 'resource', '资源质量问题已确认，要求用户重新上传', 5, 2, 'rejected'),
                                                                                                                     (15, 9, 'course_rating', '用户评分可能存在恶意刷分行为', 4, 2, 'pending'),
                                                                                                                     (18, 10, 'resource_rating', '评分与资源质量明显不符', 3, 2, 'pending');

-- 追加：资源评论的待审核记录（用于测试 /api/admin/audit/resource_comments）
INSERT INTO `reviews` (`target_id`, `reporter_id`, `target_type`, `reason`, `priority`, `reviewer_id`, `status`) VALUES
                                                                                                                     (11, 3, 'comment', '资源评论涉嫌广告，待审核', 3, NULL, 'pending'),
                                                                                                                     (16, 4, 'comment', '资源评论可能包含不当言辞，待审核', 4, NULL, 'pending'),
                                                                                                                     (20, 5, 'comment', '资源评论内容需人工复核', 3, NULL, 'pending');

-- ----------------------------
-- 16. 物品 items
-- ----------------------------
INSERT INTO `items` (`name`, `type`, `price`, `description`, `image_url`) VALUES
                                                                              ('彩虹头像框', 'avatar_frame', 100, '七彩色动态边框，彰显个性', 'https://img.example.com/frame-rainbow.png'),
                                                                              ('金色头像框', 'avatar_frame', 150, '鎏金边框，象征荣誉', 'https://img.example.com/frame-gold.png'),
                                                                              ('极简头像框', 'avatar_frame', 50, '简约白色边框，低调优雅', 'https://img.example.com/frame-minimal.png'),
                                                                              ('暗夜主题', 'background', 200, '深色系统背景，护眼模式', 'https://img.example.com/bg-dark.png'),
                                                                              ('星空主题', 'background', 300, '浩瀚星空背景，激发想象力', 'https://img.example.com/bg-starry.png'),
                                                                              ('樱花主题', 'background', 250, '春日樱花飘落背景', 'https://img.example.com/bg-sakura.png'),
                                                                              ('学霸勋章', 'medal', 500, '学习成就卓越象征', 'https://img.example.com/medal-scholar.png'),
                                                                              ('分享达人勋章', 'medal', 400, '分享资源超过10个', 'https://img.example.com/medal-sharer.png'),
                                                                              ('评论大师勋章', 'medal', 300, '发表优质评论超过50条', 'https://img.example.com/medal-commentator.png'),
                                                                              ('人气王者勋章', 'medal', 600, '获得点赞超过1000个', 'https://img.example.com/medal-popular.png'),
                                                                              ('火焰红昵称', 'nickname_color', 80, '热情如火的红色昵称', 'https://img.example.com/color-fire.png'),
                                                                              ('海洋蓝昵称', 'nickname_color', 80, '深邃如海的蓝色昵称', 'https://img.example.com/color-ocean.png'),
                                                                              ('森林绿昵称', 'nickname_color', 80, '生机勃勃的绿色昵称', 'https://img.example.com/color-forest.png');

-- ----------------------------
-- 17. 用户物品 user_items
--    使用 1/0 代替 TRUE/FALSE
-- ----------------------------
INSERT INTO `user_items` (`user_id`, `item_id`, `is_used`) VALUES
                                                               (3, 1, 1),   -- alice: 彩虹头像框（使用中）
                                                               (3, 4, 0),   -- alice: 暗夜主题
                                                               (3, 7, 1),   -- alice: 学霸勋章（使用中）
                                                               (3, 10, 0),  -- alice: 火焰红昵称
                                                               (4, 2, 0),   -- bob: 金色头像框
                                                               (4, 5, 1),   -- bob: 星空主题（使用中）
                                                               (4, 8, 0),   -- bob: 分享达人勋章
                                                               (4, 12, 1),  -- bob: 海洋蓝昵称（使用中, item_id:12）
                                                               (5, 1, 0),   -- carol: 彩虹头像框
                                                               (5, 6, 0),   -- carol: 樱花主题
                                                               (5, 9, 0),   -- carol: 评论大师勋章
                                                               (6, 3, 1),   -- david: 极简头像框（使用中）
                                                               (6, 7, 0),   -- david: 学霸勋章
                                                               (6, 13, 0),  -- david: 森林绿昵称 (item_id:13)
                                                               (7, 2, 1),   -- emma: 金色头像框（使用中）
                                                               (7, 4, 0),   -- emma: 暗夜主题
                                                               (7, 11, 0),  -- emma: 火焰红昵称 (item_id:11)
                                                               (8, 1, 0),   -- frank: 彩虹头像框
                                                               (8, 8, 1),   -- frank: 分享达人勋章（使用中）
                                                               (8, 12, 0);  -- frank: 海洋蓝昵称 (item_id:12)

-- 恢复外键检查
SET FOREIGN_KEY_CHECKS = 1;

-- ==========================================================
-- 测试数据生成完毕！
-- ==========================================================
