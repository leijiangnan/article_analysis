package handler

import (
	"article-analysis/internal/model"
	"article-analysis/internal/service"
	"net/http"
	"regexp"
	"strconv"
	"strings"
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
func (h *AnalysisHandler) AnalyzeArticle(c *gin.Context) { // ignore_security_alert
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

	// 格式化分析结果字段
	if result != nil {
		result.CoreViewpoints = formatAnalysisText(result.CoreViewpoints)
		result.FileStructure = formatAnalysisText(result.FileStructure)
		result.AuthorThoughts = formatAnalysisText(result.AuthorThoughts)
		result.RelatedMaterials = formatAnalysisText(result.RelatedMaterials)
	}

	c.JSON(http.StatusOK, model.ApiResponse{
		Code:      200,
		Message:   "success",
		Data:      result,
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
		Code:      200,
		Message:   "success",
		Data:      status,
		Timestamp: time.Now().Unix(),
	})
}

// formatAnalysisText 格式化分析文本，将数字列表格式转换为分行展示
func formatAnalysisText(text string) string {
	if text == "" {
		return text
	}

	// 如果文本中包含数字列表格式
	if strings.Contains(text, "1.") || strings.Contains(text, "1、") {
		result := text

		// 首先处理数字列表：在数字前添加换行（除了第一个数字）
		// 匹配空格或分号后的数字列表项
		result = regexp.MustCompile(`\s+(\d+)[.、]`).ReplaceAllString(result, "\n$1.")

		// 然后处理分号分隔的情况
		result = regexp.MustCompile(`(\d+)[.、]([^；;]+)[;；]\s*(\d+)[.、]`).ReplaceAllString(result, "$1.$2\n$3.")

		// 处理剩余的分号
		result = strings.ReplaceAll(result, "；", "\n")
		result = strings.ReplaceAll(result, ";", "\n")

		// 清理多余的空格和空行
		result = strings.TrimSpace(result)
		// 移除连续的空行
		result = regexp.MustCompile(`\n\s*\n`).ReplaceAllString(result, "\n")

		return result
	}

	return text
}
