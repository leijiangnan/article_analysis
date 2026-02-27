package handler

import (
	"article-analysis/internal/model"
	"article-analysis/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CrawlerHandler struct {
	crawlerService *service.CrawlerService
}

func NewCrawlerHandler(crawlerService *service.CrawlerService) *CrawlerHandler {
	return &CrawlerHandler{
		crawlerService: crawlerService,
	}
}

type CrawlRequest struct {
	URL   string `json:"url" binding:"required,url"`
	Count int    `json:"count" binding:"required,min=1,max=100"`
}

type CrawlResponse struct {
	TotalFound   int                     `json:"total_found"`
	CrawledCount int                     `json:"crawled_count"`
	SavedCount   int                     `json:"saved_count"`
	Articles     []CrawledArticleSummary `json:"articles"`
	Errors       []string                `json:"errors,omitempty"`
}

type CrawledArticleSummary struct {
	ID     uint64 `json:"id,omitempty"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Date   string `json:"date,omitempty"`
	URL    string `json:"url"`
	Saved  bool   `json:"saved"`
}

func (h *CrawlerHandler) CrawlArticles(c *gin.Context) {
	var req CrawlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ApiResponse{
			Code:      400,
			Message:   "参数错误: " + err.Error(),
			Timestamp: time.Now().Unix(),
		})
		return
	}

	result, err := h.crawlerService.CrawlArticles(req.URL, req.Count)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ApiResponse{
			Code:      400,
			Message:   err.Error(),
			Timestamp: time.Now().Unix(),
		})
		return
	}

	savedArticles, saveErrors := h.crawlerService.SaveCrawledArticles(result.Articles)

	allErrors := append(result.Errors, make([]string, 0)...)
	for _, e := range saveErrors {
		allErrors = append(allErrors, e.Error())
	}

	articles := make([]CrawledArticleSummary, 0, len(result.Articles))
	savedMap := make(map[string]uint64)
	for _, a := range savedArticles {
		savedMap[a.Title] = a.ID
	}

	for _, a := range result.Articles {
		summary := CrawledArticleSummary{
			Title:  a.Title,
			Author: a.Author,
			Date:   a.Date,
			URL:    a.URL,
		}
		if id, ok := savedMap[a.Title]; ok {
			summary.ID = id
			summary.Saved = true
		}
		articles = append(articles, summary)
	}

	response := CrawlResponse{
		TotalFound:   result.TotalFound,
		CrawledCount: result.CrawledCount,
		SavedCount:   len(savedArticles),
		Articles:     articles,
		Errors:       allErrors,
	}

	c.JSON(http.StatusOK, model.ApiResponse{
		Code:      200,
		Message:   "爬取完成",
		Data:      response,
		Timestamp: time.Now().Unix(),
	})
}
