package handler

import (
	"article-analysis/internal/model"
	"article-analysis/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AnalysisHandler struct {
	analysisService *service.AnalysisService
}

func NewAnalysisHandler(analysisService *service.AnalysisService) *AnalysisHandler {
	return &AnalysisHandler{
		analysisService: analysisService,
	}
}

// AnalyzeArticle 提交文章分析
func (h *AnalysisHandler) AnalyzeArticle(c *gin.Context) {
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
	
	task, err := h.analysisService.AnalyzeArticle(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ApiResponse{
			Code:      422,
			Message:   err.Error(),
			Timestamp: time.Now().Unix(),
		})
		return
	}
	
	c.JSON(http.StatusOK, model.ApiResponse{
		Code:    200,
		Message: "分析任务已提交",
		Data: map[string]interface{}{
			"task_id": task.TaskID,
			"status":  task.Status,
		},
		Timestamp: time.Now().Unix(),
	})
}

// GetAnalysisResult 获取分析结果
func (h *AnalysisHandler) GetAnalysisResult(c *gin.Context) {
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
	
	result, err := h.analysisService.GetAnalysisResult(id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.ApiResponse{
			Code:      404,
			Message:   err.Error(),
			Timestamp: time.Now().Unix(),
		})
		return
	}
	
	c.JSON(http.StatusOK, model.ApiResponse{
		Code:    200,
		Message: "success",
		Data:    result,
		Timestamp: time.Now().Unix(),
	})
}

// GetAnalysisStatus 获取分析任务状态
func (h *AnalysisHandler) GetAnalysisStatus(c *gin.Context) {
	taskID := c.Param("task_id")
	
	status, err := h.analysisService.GetAnalysisStatus(taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, model.ApiResponse{
			Code:      404,
			Message:   "任务不存在",
			Timestamp: time.Now().Unix(),
		})
		return
	}
	
	c.JSON(http.StatusOK, model.ApiResponse{
		Code:    200,
		Message: "success",
		Data:    status,
		Timestamp: time.Now().Unix(),
	})
}