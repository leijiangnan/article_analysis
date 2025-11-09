package service

import (
	"context"
	"encoding/json"
	"testing"

	"article-analysis/internal/config"
	"article-analysis/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOpenAIClient 是一个模拟的OpenAI客户端
type MockOpenAIClient struct {
	mock.Mock
}

func (m *MockOpenAIClient) AnalyzeArticle(ctx context.Context, content string) (*AnalysisResponse, error) {
	args := m.Called(ctx, content)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AnalysisResponse), args.Error(1)
}

func TestOpenAIClient_AnalyzeArticle_Success(t *testing.T) {
	// 创建mock logger
	log := logger.NewLogger("test")
	
	// 创建OpenAI客户端
	cfg := &config.Config{
		OpenAI: config.OpenAIConfig{
			APIKey: "test-api-key",
		},
	}
	
	client := NewOpenAIClient(cfg, log)
	assert.NotNil(t, client)
}

func TestOpenAIClient_buildAnalysisPrompt(t *testing.T) {
	log := logger.NewLogger("test")
	cfg := &config.Config{
		OpenAI: config.OpenAIConfig{
			APIKey: "test-api-key",
		},
	}
	
	client := NewOpenAIClient(cfg, log)
	
	content := "这是一篇测试文章内容"
	prompt := client.buildAnalysisPrompt(content)
	
	// 验证提示词包含必要的内容
	assert.Contains(t, prompt, "请对以下文章进行深度分析")
	assert.Contains(t, prompt, content)
	assert.Contains(t, prompt, "核心观点")
	assert.Contains(t, prompt, "文件结构")
	assert.Contains(t, prompt, "作者思路")
	assert.Contains(t, prompt, "相关素材与事例")
	assert.Contains(t, prompt, "JSON格式返回结果")
}

func TestOpenAIClient_parseAIResponse_Success(t *testing.T) {
	log := logger.NewLogger("test")
	cfg := &config.Config{
		OpenAI: config.OpenAIConfig{
			APIKey: "test-api-key",
		},
	}
	
	client := NewOpenAIClient(cfg, log)
	
	// 测试有效的JSON响应
	responseContent := `
以下是分析结果：
{
  "core_viewpoints": "文章核心观点",
  "file_structure": "文章结构描述",
  "author_thoughts": "作者思路分析",
  "related_materials": "相关素材"
}
希望这个分析对您有帮助！`
	
	result, err := client.parseAIResponse(responseContent)
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "文章核心观点", result.CoreViewpoints)
	assert.Equal(t, "文章结构描述", result.FileStructure)
	assert.Equal(t, "作者思路分析", result.AuthorThoughts)
	assert.Equal(t, "相关素材", result.RelatedMaterials)
}

func TestOpenAIClient_parseAIResponse_InvalidJSON(t *testing.T) {
	log := logger.NewLogger("test")
	cfg := &config.Config{
		OpenAI: config.OpenAIConfig{
			APIKey: "test-api-key",
		},
	}
	
	client := NewOpenAIClient(cfg, log)
	
	// 测试无效的JSON响应
	responseContent := `
以下是分析结果：
{
  "invalid_json": "缺少闭合括号"
希望这个分析对您有帮助！`
	
	result, err := client.parseAIResponse(responseContent)
	
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "无法解析AI响应格式")
}

func TestOpenAIClient_parseAIResponse_NoJSON(t *testing.T) {
	log := logger.NewLogger("test")
	cfg := &config.Config{
		OpenAI: config.OpenAIConfig{
			APIKey: "test-api-key",
		},
	}
	
	client := NewOpenAIClient(cfg, log)
	
	// 测试没有JSON的响应
	responseContent := `
以下是分析结果：
没有JSON格式的内容
希望这个分析对您有帮助！`
	
	result, err := client.parseAIResponse(responseContent)
	
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "无法解析AI响应格式")
}

func TestOpenAIClient_getModel(t *testing.T) {
	log := logger.NewLogger("test")
	
	// 测试配置了模型的情况
	cfg1 := &config.Config{
		OpenAI: config.OpenAIConfig{
			APIKey: "test-api-key",
			Model:  "gpt-4",
		},
	}
	client1 := NewOpenAIClient(cfg1, log)
	assert.Equal(t, "gpt-4", client1.getModel())
	
	// 测试未配置模型的情况（使用默认）
	cfg2 := &config.Config{
		OpenAI: config.OpenAIConfig{
			APIKey: "test-api-key",
			Model:  "",
		},
	}
	client2 := NewOpenAIClient(cfg2, log)
	assert.Equal(t, "kimi-k2-0905-preview", client2.getModel())
}

func TestAnalysisResult_JSONMarshal(t *testing.T) {
	result := &AnalysisResponse{
		CoreViewpoints:   "核心观点",
		FileStructure:    "文件结构",
		AuthorThoughts:   "作者思路",
		RelatedMaterials: "相关素材",
	}
	
	jsonData, err := json.Marshal(result)
	assert.NoError(t, err)
	
	// 验证JSON格式
	var unmarshaled AnalysisResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, result.CoreViewpoints, unmarshaled.CoreViewpoints)
	assert.Equal(t, result.FileStructure, unmarshaled.FileStructure)
	assert.Equal(t, result.AuthorThoughts, unmarshaled.AuthorThoughts)
	assert.Equal(t, result.RelatedMaterials, unmarshaled.RelatedMaterials)
}