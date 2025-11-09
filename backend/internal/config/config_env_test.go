package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_EnvVariablePriority(t *testing.T) {
	// 保存原始环境变量
	originalAPIKey := os.Getenv("OPENAI_API_KEY")
	originalAPIBase := os.Getenv("OPENAI_API_BASE")
	
	// 测试完成后恢复原始环境变量
	defer func() {
		os.Setenv("OPENAI_API_KEY", originalAPIKey)
		os.Setenv("OPENAI_API_BASE", originalAPIBase)
	}()
	
	// 设置测试环境变量
	os.Setenv("OPENAI_API_KEY", "env-api-key-123")
	os.Setenv("OPENAI_API_BASE", "https://env-api.openai.com/v1")
	
	// 加载配置
	config, err := LoadConfig()
	
	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, config)
	
	// 验证环境变量优先级
	assert.Equal(t, "env-api-key-123", config.OpenAI.APIKey, "应该使用环境变量中的API密钥")
	assert.Equal(t, "https://env-api.openai.com/v1", config.OpenAI.APIBase, "应该使用环境变量中的API基础URL")
}

func TestLoadConfig_NoEnvVariable(t *testing.T) {
	// 保存原始环境变量
	originalAPIKey := os.Getenv("OPENAI_API_KEY")
	originalAPIBase := os.Getenv("OPENAI_API_BASE")
	
	// 测试完成后恢复原始环境变量
	defer func() {
		os.Setenv("OPENAI_API_KEY", originalAPIKey)
		os.Setenv("OPENAI_API_BASE", originalAPIBase)
	}()
	
	// 清空环境变量
	os.Unsetenv("OPENAI_API_KEY")
	os.Unsetenv("OPENAI_API_BASE")
	
	// 加载配置
	config, err := LoadConfig()
	
	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, config)
	
	// 验证使用默认值
	assert.Equal(t, "", config.OpenAI.APIKey, "没有环境变量时API密钥应该为空")
	assert.Equal(t, "https://api.moonshot.cn/v1", config.OpenAI.APIBase, "没有环境变量时应该使用默认API基础URL")
}