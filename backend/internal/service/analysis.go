package service

import (
	"article-analysis/internal/config"
	"article-analysis/internal/model"
	"article-analysis/internal/repository"
	"article-analysis/pkg/logger"
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type AnalysisService struct {
	analysisRepo *repository.AnalysisRepository
	articleRepo  *repository.ArticleRepository
	openaiClient *OpenAIClient
	log          *logger.Logger
}

func NewAnalysisService(analysisRepo *repository.AnalysisRepository, articleRepo *repository.ArticleRepository, cfg *config.Config, log *logger.Logger) *AnalysisService {
	return &AnalysisService{
		analysisRepo: analysisRepo,
		articleRepo:  articleRepo,
		openaiClient: NewOpenAIClient(cfg, log),
		log:          log,
	}
}

type AnalysisTask struct {
	TaskID    string
	ArticleID uint64
	Status    string
}

func (s *AnalysisService) AnalyzeArticle(articleID uint64) (*AnalysisTask, error) {
	// 检查文章是否存在
	article, err := s.articleRepo.GetByID(articleID)
	if err != nil {
		return nil, errors.New("文章不存在")
	}
	
	// 检查是否已有分析任务
	existingAnalysis, err := s.analysisRepo.GetByArticleID(articleID)
	if err == nil && existingAnalysis.AnalysisStatus == "processing" {
		return nil, errors.New("分析任务正在进行中")
	}
	
	// 创建分析任务
	taskID := fmt.Sprintf("task_%d_%d", articleID, time.Now().Unix())
	analysis := &model.ArticleAnalysis{
		ArticleID:      articleID,
		AnalysisStatus: "processing", // 直接设置为处理中
	}
	
	if existingAnalysis == nil {
		if err := s.analysisRepo.Create(analysis); err != nil {
			return nil, errors.New("创建分析任务失败")
		}
	} else {
		analysis = existingAnalysis
		analysis.AnalysisStatus = "processing"
		analysis.ErrorMessage = ""
		if err := s.analysisRepo.Update(analysis); err != nil {
			return nil, errors.New("更新分析任务失败")
		}
	}
	
	// 异步执行分析任务
	go s.performAnalysis(articleID, article.Content)
	
	return &AnalysisTask{
		TaskID:    taskID,
		ArticleID: articleID,
		Status:    "processing",
	}, nil
}

func (s *AnalysisService) performAnalysis(articleID uint64, content string) {
	// 更新状态为处理中
	if err := s.analysisRepo.UpdateStatus(articleID, "processing", ""); err != nil {
		s.log.Error("更新分析状态失败", err)
		return
	}
	
	// 使用OpenAI客户端进行分析
	analysisResult, err := s.openaiClient.AnalyzeArticle(timeoutContext(), content)
	if err != nil {
		s.log.Error("AI分析失败", err)
		s.analysisRepo.UpdateStatus(articleID, "failed", fmt.Sprintf("AI分析失败: %v", err))
		return
	}
	
	// 保存分析结果
	analysis, err := s.analysisRepo.GetByArticleID(articleID)
	if err != nil {
		s.log.Error("获取分析记录失败", err)
		return
	}
	
	analysis.CoreViewpoints = analysisResult.CoreViewpoints
	analysis.FileStructure = analysisResult.FileStructure
	analysis.AuthorThoughts = analysisResult.AuthorThoughts
	analysis.RelatedMaterials = analysisResult.RelatedMaterials
	analysis.AnalysisStatus = "completed"
	analysis.ErrorMessage = ""
	
	if err := s.analysisRepo.Update(analysis); err != nil {
		s.log.Error("保存分析结果失败", err)
		s.analysisRepo.UpdateStatus(articleID, "failed", "保存分析结果失败")
		return
	}
	
	s.log.Info("文章分析完成", zap.Int("article_id", int(articleID)))
}

func (s *AnalysisService) GetAnalysisResult(articleID uint64) (*model.ArticleAnalysis, error) {
	analysis, err := s.analysisRepo.GetByArticleID(articleID)
	if err != nil {
		return nil, errors.New("分析结果不存在")
	}
	return analysis, nil
}

func (s *AnalysisService) GetAnalysisStatus(taskID string) (map[string]interface{}, error) {
	// 从任务ID中提取文章ID
	var articleID uint64
	fmt.Sscanf(taskID, "task_%d_", &articleID)
	
	if articleID == 0 {
		return nil, errors.New("无效的任务ID")
	}
	
	// 查询数据库获取实际状态
	analysis, err := s.analysisRepo.GetByArticleID(articleID)
	if err != nil {
		return map[string]interface{}{
			"task_id":  taskID,
			"status":   "pending",
			"progress": 0,
		}, nil
	}
	
	progress := 0
	switch analysis.AnalysisStatus {
	case "completed":
		progress = 100
	case "processing":
		progress = 50
	case "failed":
		progress = 0
	default:
		progress = 0
	}
	
	return map[string]interface{}{
		"task_id":  taskID,
		"status":   analysis.AnalysisStatus,
		"progress": progress,
		"error":    analysis.ErrorMessage,
	}, nil
}

func timeoutContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 120*time.Second)
	return ctx
}