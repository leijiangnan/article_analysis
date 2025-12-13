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
	config *config.Config
}

func NewOpenAIClient(cfg *config.Config, log *logger.Logger) *OpenAIClient {
	clientConfig := openai.DefaultConfig(cfg.OpenAI.APIKey)
	if cfg.OpenAI.APIBase != "" {
		clientConfig.BaseURL = cfg.OpenAI.APIBase
	}

	return &OpenAIClient{
		client: openai.NewClientWithConfig(clientConfig),
		log:    log,
		config: cfg,
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

	// 使用配置的模型，如果没有配置则使用默认的GPT-3.5-turbo
	model := c.getModel()

	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: model,
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
你是一个专业的文章分析助手。

任务：对“文章内容”进行深度分析，并输出分析结果。

输出要求（必须严格遵守）：
1) 只输出一个合法的JSON对象，输出必须以“{”开头、以“}”结尾；不要输出任何额外文本、解释、Markdown、代码块标记（例如不要输出“三个反引号”）。
2) JSON只允许包含以下四个key（键名必须完全一致、区分大小写），value必须是字符串：
   - "core_viewpoints"
   - "file_structure"
   - "author_thoughts"
   - "related_materials"
3) 四个字段都必须给出内容；如果原文信息不足，请基于原文做合理推断并说明“不确定点”，不要留空。
4) 如果某个字段需要列点说明：在字符串中使用“1. … 2. … 3. …”这样的编号格式。
5) 全部使用中文表达。

文章内容：
%s

请输出严格符合以下JSON结构的结果：
{
  "core_viewpoints": "...",
  "file_structure": "...",
  "author_thoughts": "...",
  "related_materials": "..."
}

（参考：文章内容长度 %d 字符）
`, content, len(content))
}

func (c *OpenAIClient) getModel() string {
	if c.config.OpenAI.Model != "" {
		return c.config.OpenAI.Model
	}
	return "kimi-k2-0905-preview" // 默认Moonshot模型
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
