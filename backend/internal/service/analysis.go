package service

import (
	"article-analysis/internal/config"
	"article-analysis/internal/model"
	"article-analysis/internal/repository"
	"article-analysis/pkg/logger"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

type AnalysisService struct {
	analysisRepo *repository.AnalysisRepository
	articleRepo  *repository.ArticleRepository
	openaiClient *openai.Client
	log          *logger.Logger
}

func NewAnalysisService(analysisRepo *repository.AnalysisRepository, articleRepo *repository.ArticleRepository, cfg *config.Config, log *logger.Logger) *AnalysisService {
	client := openai.NewClient(cfg.OpenAI.APIKey)
	return &AnalysisService{
		analysisRepo: analysisRepo,
		articleRepo:  articleRepo,
		openaiClient: client,
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
		AnalysisStatus: "pending",
	}
	
	if existingAnalysis == nil {
		if err := s.analysisRepo.Create(analysis); err != nil {
			return nil, errors.New("创建分析任务失败")
		}
	} else {
		analysis = existingAnalysis
		analysis.AnalysisStatus = "pending"
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
	
	// 构建分析提示词
	prompt := s.buildAnalysisPrompt(content)
	
	// 调用OpenAI API
	resp, err := s.openaiClient.CreateChatCompletion(
		timeoutContext(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "你是一个专业的文章分析助手，请对文章内容进行深度分析。",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens: 2000,
			Temperature: 0.7,
		},
	)
	
	if err != nil {
		s.log.Error("OpenAI API调用失败", err)
		s.analysisRepo.UpdateStatus(articleID, "failed", fmt.Sprintf("AI分析失败: %v", err))
		return
	}
	
	// 解析AI响应
	analysisResult, err := s.parseAIResponse(resp.Choices[0].Message.Content)
	if err != nil {
		s.log.Error("解析AI响应失败", err)
		s.analysisRepo.UpdateStatus(articleID, "failed", fmt.Sprintf("解析分析结果失败: %v", err))
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
	
	if err := s.analysisRepo.Update(analysis); err != nil {
		s.log.Error("保存分析结果失败", err)
		s.analysisRepo.UpdateStatus(articleID, "failed", "保存分析结果失败")
		return
	}
	
	s.log.Info("文章分析完成", zap.Int("article_id", int(articleID)))
}

func (s *AnalysisService) buildAnalysisPrompt(content string) string {
	return fmt.Sprintf(`
请对以下文章进行深度分析，并以JSON格式返回分析结果：

文章内容：
%s

请提供以下四个方面的分析：

1. 核心观点：总结文章的主要观点和核心论点
2. 文件结构：分析文章的结构组织方式
3. 作者思路：分析作者的写作思路和逻辑脉络
4. 相关素材与事例：提取文章中的重要素材、案例和论据

请以以下JSON格式返回结果：
{
  "core_viewpoints": "核心观点内容",
  "file_structure": "文件结构描述",
  "author_thoughts": "作者思路分析",
  "related_materials": "相关素材与事例"
}

文章内容长度：%d字符`, content, len(content))
}

type AnalysisResult struct {
	CoreViewpoints   string `json:"core_viewpoints"`
	FileStructure    string `json:"file_structure"`
	AuthorThoughts   string `json:"author_thoughts"`
	RelatedMaterials string `json:"related_materials"`
}

func (s *AnalysisService) parseAIResponse(content string) (*AnalysisResult, error) {
	// 提取JSON部分
	startIdx := strings.Index(content, "{")
	endIdx := strings.LastIndex(content, "}")
	
	if startIdx == -1 || endIdx == -1 {
		return nil, errors.New("无法解析AI响应格式")
	}
	
	jsonStr := content[startIdx : endIdx+1]
	
	var result AnalysisResult
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %w", err)
	}
	
	return &result, nil
}

func (s *AnalysisService) GetAnalysisResult(articleID uint64) (*model.ArticleAnalysis, error) {
	analysis, err := s.analysisRepo.GetByArticleID(articleID)
	if err != nil {
		return nil, errors.New("分析结果不存在")
	}
	return analysis, nil
}

func (s *AnalysisService) GetAnalysisStatus(taskID string) (map[string]interface{}, error) {
	// 简单的状态查询，实际项目中可以使用Redis等存储任务状态
	return map[string]interface{}{
		"task_id": taskID,
		"status":  "completed", // 简化处理
		"progress": 100,
	}, nil
}

func timeoutContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	return ctx
}