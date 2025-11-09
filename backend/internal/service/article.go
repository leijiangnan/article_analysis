package service

import (
	"article-analysis/internal/model"
	"article-analysis/internal/repository"
	"article-analysis/pkg/logger"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
)

type ArticleService struct {
	repo *repository.ArticleRepository
	log  *logger.Logger
}

func NewArticleService(repo *repository.ArticleRepository, log *logger.Logger) *ArticleService {
	return &ArticleService{
		repo: repo,
		log:  log,
	}
}

func (s *ArticleService) UploadArticle(file *multipart.FileHeader, title, author string) (*model.Article, error) {
	// 验证文件类型
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".txt") {
		return nil, errors.New("只支持TXT格式文件")
	}
	
	// 限制文件大小 (10MB)
	if file.Size > 10*1024*1024 {
		return nil, errors.New("文件大小不能超过10MB")
	}
	
	// 读取文件内容
	src, err := file.Open()
	if err != nil {
		s.log.Error("打开文件失败", err)
		return nil, errors.New("文件读取失败")
	}
	defer src.Close()
	
	content, err := io.ReadAll(src)
	if err != nil {
		s.log.Error("读取文件内容失败", err)
		return nil, errors.New("文件内容读取失败")
	}
	
	// 自动提取标题和作者（如果未提供）
	if title == "" {
		title = s.extractTitleFromContent(string(content), file.Filename)
	}
	if author == "" {
		author = s.extractAuthorFromContent(string(content))
	}
	
	// 保存文件到本地
	uploadDir := "./web/uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		s.log.Error("创建上传目录失败", err)
		return nil, errors.New("文件保存失败")
	}
	
	filename := time.Now().Format("20060102_150405_") + file.Filename
	filePath := filepath.Join(uploadDir, filename)
	
	// 重新打开文件进行保存
	src.Seek(0, 0)
	dst, err := os.Create(filePath)
	if err != nil {
		s.log.Error("创建文件失败", err)
		return nil, errors.New("文件保存失败")
	}
	defer dst.Close()
	
	if _, err := io.Copy(dst, src); err != nil {
		s.log.Error("保存文件失败", err)
		return nil, errors.New("文件保存失败")
	}
	
	// 创建文章记录
	article := &model.Article{
		Title:    title,
		Author:   author,
		Content:  string(content),
		FilePath: filePath,
		FileSize: file.Size,
	}
	
	if err := s.repo.Create(article); err != nil {
		s.log.Error("保存文章记录失败", err)
		// 清理已保存的文件
		os.Remove(filePath)
		return nil, errors.New("文章保存失败")
	}
	
	s.log.Info("文章上传成功", 
		zap.String("title", title), 
		zap.String("author", author),
		zap.Int("size", int(file.Size)))
	
	return article, nil
}

func (s *ArticleService) GetArticleList(req *model.PaginationRequest) (*model.PaginationResponse, error) {
	return s.repo.GetList(req)
}

func (s *ArticleService) GetArticleDetail(id uint64) (*model.Article, error) {
	article, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("文章不存在")
	}
	return article, nil
}

func (s *ArticleService) GetAuthors() ([]string, error) {
	return s.repo.GetAuthors()
}

// DeleteArticle 删除文章
func (s *ArticleService) DeleteArticle(id uint64) error {
	// 首先获取文章信息
	article, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("文章不存在")
	}
	
	// 删除关联的文件
	if article.FilePath != "" {
		if err := os.Remove(article.FilePath); err != nil {
			s.log.Warn("删除文件失败", zap.String("file_path", article.FilePath), zap.Error(err))
			// 继续删除数据库记录，即使文件删除失败
		}
	}
	
	// 删除数据库记录
	if err := s.repo.Delete(id); err != nil {
		s.log.Error("删除文章数据库记录失败", err)
		return errors.New("删除文章失败")
	}
	
	s.log.Info("文章删除成功", zap.Uint64("id", id), zap.String("title", article.Title))
	return nil
}

// 从内容中提取标题
func (s *ArticleService) extractTitleFromContent(content, filename string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 10 && len(line) < 100 {
			return line
		}
	}
	// 使用文件名作为标题
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

// 从内容中提取作者
func (s *ArticleService) extractAuthorFromContent(content string) string {
	// 简单的作者识别逻辑
	authorKeywords := []string{"作者：", "作者:", "Author:", "By:"}
	lines := strings.Split(content, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		for _, keyword := range authorKeywords {
			if strings.Contains(line, keyword) {
				author := strings.TrimSpace(strings.Replace(line, keyword, "", 1))
				if len(author) > 0 && len(author) < 50 {
					return author
				}
			}
		}
	}
	
	return "未知作者"
}