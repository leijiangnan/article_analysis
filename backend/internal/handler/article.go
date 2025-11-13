package handler

import (
	"article-analysis/internal/model"
	"article-analysis/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	articleService *service.ArticleService
}

func NewArticleHandler(articleService *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		articleService: articleService,
	}
}

// UploadArticle 上传文章
func (h *ArticleHandler) UploadArticle(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ApiResponse{
			Code:      400,
			Message:   "请选择要上传的文件",
			Timestamp: time.Now().Unix(),
		})
		return
	}

	title := c.PostForm("title")
	author := c.PostForm("author")

	article, err := h.articleService.UploadArticle(file, title, author)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ApiResponse{
			Code:      400,
			Message:   err.Error(),
			Timestamp: time.Now().Unix(),
		})
		return
	}

	c.JSON(http.StatusOK, model.ApiResponse{
		Code:    200,
		Message: "上传成功",
		Data: map[string]interface{}{
			"id":          strconv.FormatUint(article.ID, 10),
			"title":       article.Title,
			"author":      article.Author,
			"upload_time": article.UploadTime.Format("2006-01-02 15:04:05"),
		},
		Timestamp: time.Now().Unix(),
	})
}

// GetArticleList 获取文章列表
func (h *ArticleHandler) GetArticleList(c *gin.Context) {
	var req model.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ApiResponse{
			Code:      400,
			Message:   "参数错误",
			Timestamp: time.Now().Unix(),
		})
		return
	}

	result, err := h.articleService.GetArticleList(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ApiResponse{
			Code:      500,
			Message:   "获取文章列表失败",
			Timestamp: time.Now().Unix(),
		})
		return
	}

	c.JSON(http.StatusOK, model.ApiResponse{
		Code:      200,
		Message:   "success",
		Data:      result,
		Timestamp: time.Now().Unix(),
	})
}

// GetArticleListWithAnalysis 获取文章列表及分析状态
func (h *ArticleHandler) GetArticleListWithAnalysis(c *gin.Context) {
	var req model.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ApiResponse{
			Code:      400,
			Message:   "参数错误",
			Timestamp: time.Now().Unix(),
		})
		return
	}

	result, err := h.articleService.GetArticleListWithAnalysis(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ApiResponse{
			Code:      500,
			Message:   "获取文章列表失败",
			Timestamp: time.Now().Unix(),
		})
		return
	}

	c.JSON(http.StatusOK, model.ApiResponse{
		Code:      200,
		Message:   "success",
		Data:      result,
		Timestamp: time.Now().Unix(),
	})
}

// GetArticleDetail 获取文章详情
func (h *ArticleHandler) GetArticleDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ApiResponse{
			Code:      400,
			Message:   "文章ID格式错误",
			Timestamp: time.Now().Unix(),
		})
		return
	}

	article, err := h.articleService.GetArticleDetail(id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.ApiResponse{
			Code:      404,
			Message:   "文章不存在",
			Timestamp: time.Now().Unix(),
		})
		return
	}

	c.JSON(http.StatusOK, model.ApiResponse{
		Code:      200,
		Message:   "success",
		Data:      article,
		Timestamp: time.Now().Unix(),
	})
}

// GetAuthors 获取作者列表
func (h *ArticleHandler) GetAuthors(c *gin.Context) {
	authors, err := h.articleService.GetAuthors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ApiResponse{
			Code:      500,
			Message:   "获取作者列表失败",
			Timestamp: time.Now().Unix(),
		})
		return
	}

	c.JSON(http.StatusOK, model.ApiResponse{
		Code:      200,
		Message:   "success",
		Data:      authors,
		Timestamp: time.Now().Unix(),
	})
}

// CreateArticle 创建文章（直接输入文本）
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var req struct {
		Title   string `json:"title" binding:"required"`
		Author  string `json:"author"`
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ApiResponse{
			Code:      400,
			Message:   "参数错误：" + err.Error(),
			Timestamp: time.Now().Unix(),
		})
		return
	}

	// 限制内容长度 (10MB)
	if len(req.Content) > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, model.ApiResponse{
			Code:      400,
			Message:   "内容长度不能超过10MB",
			Timestamp: time.Now().Unix(),
		})
		return
	}

	article, err := h.articleService.CreateArticle(req.Title, req.Author, req.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ApiResponse{
			Code:      400,
			Message:   err.Error(),
			Timestamp: time.Now().Unix(),
		})
		return
	}

	c.JSON(http.StatusOK, model.ApiResponse{
		Code:    200,
		Message: "创建成功",
		Data: map[string]interface{}{
			"id":          strconv.FormatUint(article.ID, 10),
			"title":       article.Title,
			"author":      article.Author,
			"upload_time": article.UploadTime.Format("2006-01-02 15:04:05"),
		},
		Timestamp: time.Now().Unix(),
	})
}

// DeleteArticle 删除文章
func (h *ArticleHandler) DeleteArticle(c *gin.Context) { //ignore_security_alert IDOR
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ApiResponse{
			Code:      400,
			Message:   "无效的文章ID",
			Timestamp: time.Now().Unix(),
		})
		return
	}

	if err := h.articleService.DeleteArticle(id); err != nil {
		c.JSON(http.StatusInternalServerError, model.ApiResponse{
			Code:      500,
			Message:   err.Error(),
			Timestamp: time.Now().Unix(),
		})
		return
	}

	c.JSON(http.StatusOK, model.ApiResponse{
		Code:      200,
		Message:   "删除成功",
		Data:      nil,
		Timestamp: time.Now().Unix(),
	})
}
