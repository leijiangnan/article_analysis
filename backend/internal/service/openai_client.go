package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"article-analysis/internal/config"
	"article-analysis/pkg/logger"

	"github.com/sashabaranov/go-openai"
)

type OpenAIClient struct {
	client *openai.Client
	log    *logger.Logger
}

func NewOpenAIClient(cfg *config.Config, log *logger.Logger) *OpenAIClient {
	return &OpenAIClient{
		client: openai.NewClient(cfg.OpenAI.APIKey),
		log:    log,
	}
}

type AnalysisRequest struct {
	Content string
	Prompt  string
}

type AnalysisResponse struct {
	CoreViewpoints   string
	FileStructure    string
	AuthorThoughts   string
	RelatedMaterials string
}

func (c *OpenAIClient) AnalyzeArticle(ctx context.Context, content string) (*AnalysisResponse, error) {
	prompt := c.buildAnalysisPrompt(content)
	
	resp, err := c.client.CreateChatCompletion(
		ctx,
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
		c.log.Error("OpenAI API调用失败", err)
		return nil, fmt.Errorf("OpenAI API调用失败: %w", err)
	}
	
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("OpenAI API返回空响应")
	}
	
	result, err := c.parseAIResponse(resp.Choices[0].Message.Content)
	if err != nil {
		c.log.Error("解析AI响应失败", err)
		return nil, fmt.Errorf("解析AI响应失败: %w", err)
	}
	
	return result, nil
}

func (c *OpenAIClient) buildAnalysisPrompt(content string) string {
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

func (c *OpenAIClient) parseAIResponse(content string) (*AnalysisResponse, error) {
	// 提取JSON部分
	startIdx := strings.Index(content, "{")
	endIdx := strings.LastIndex(content, "}")
	
	if startIdx == -1 || endIdx == -1 {
		return nil, fmt.Errorf("无法解析AI响应格式，找不到JSON内容")
	}
	
	jsonStr := content[startIdx : endIdx+1]
	
	// 定义临时结构体用于解析
	tempResult := struct {
		CoreViewpoints   string `json:"core_viewpoints"`
		FileStructure    string `json:"file_structure"`
		AuthorThoughts   string `json:"author_thoughts"`
		RelatedMaterials string `json:"related_materials"`
	}{}
	
	if err := json.Unmarshal([]byte(jsonStr), &tempResult); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %w", err)
	}
	
	return &AnalysisResponse{
		CoreViewpoints:   tempResult.CoreViewpoints,
		FileStructure:    tempResult.FileStructure,
		AuthorThoughts:   tempResult.AuthorThoughts,
		RelatedMaterials: tempResult.RelatedMaterials,
	}, nil
}