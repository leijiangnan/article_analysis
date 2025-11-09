-- 创建数据库
CREATE DATABASE IF NOT EXISTS article_analysis CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE article_analysis;

-- 创建文章表
CREATE TABLE IF NOT EXISTS articles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(500) NOT NULL COMMENT '文章标题',
    author VARCHAR(200) NOT NULL COMMENT '作者',
    content TEXT NOT NULL COMMENT '文章内容',
    file_path VARCHAR(500) NOT NULL COMMENT '文件存储路径',
    file_size BIGINT NOT NULL COMMENT '文件大小(字节)',
    upload_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_author (author),
    INDEX idx_upload_time (upload_time),
    FULLTEXT idx_title_content (title, content)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章表';

-- 创建文章分析结果表
CREATE TABLE IF NOT EXISTS article_analyses (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    article_id BIGINT NOT NULL COMMENT '文章ID',
    core_viewpoints TEXT COMMENT '核心观点',
    file_structure TEXT COMMENT '文件结构',
    author_thoughts TEXT COMMENT '作者思路',
    related_materials TEXT COMMENT '相关素材与事例',
    analysis_status ENUM('pending','processing','completed','failed') DEFAULT 'pending' COMMENT '分析状态',
    analysis_time TIMESTAMP NULL COMMENT '分析完成时间',
    error_message TEXT COMMENT '错误信息',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE,
    INDEX idx_article_id (article_id),
    INDEX idx_status (analysis_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章分析结果表';

-- 插入测试数据
INSERT INTO articles (title, author, content, file_path, file_size) VALUES
('人工智能的未来发展', '张三', '人工智能技术正在快速发展...', '/uploads/test1.txt', 1024),
('机器学习基础教程', '李四', '机器学习是人工智能的重要分支...', '/uploads/test2.txt', 2048),
('深度学习实践指南', '王五', '深度学习在图像识别等领域有广泛应用...', '/uploads/test3.txt', 1536);

-- 创建分析结果测试数据
INSERT INTO article_analyses (article_id, core_viewpoints, file_structure, author_thoughts, related_materials, analysis_status, analysis_time) VALUES
(1, 'AI将改变人类社会的各个方面', '总-分-总结构', '从宏观角度分析AI发展趋势', '引用了大量研究数据和案例', 'completed', NOW()),
(2, '机器学习是AI的核心技术', '教程式结构', '循序渐进地介绍ML概念', '提供了实际代码示例', 'completed', NOW());