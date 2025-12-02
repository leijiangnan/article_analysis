package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Server   ServerConfig   `mapstructure:"server"`
	OpenAI   OpenAIConfig   `mapstructure:"openai"`
	Log      LogConfig      `mapstructure:"log"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Driver   string `mapstructure:"driver"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type OpenAIConfig struct {
	APIKey  string `mapstructure:"api_key"`
	APIBase string `mapstructure:"api_base"`
	Model   string `mapstructure:"model"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
}

func LoadConfig() (*Config, error) {
	// 首先检查是否通过环境变量指定了配置文件路径
	configPath := viper.GetString("CONFIG_PATH")
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		// 使用默认的配置文件名和搜索路径
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
	}

	// 确保也能从环境变量读取配置路径
	viper.AutomaticEnv()
	viper.BindEnv("CONFIG_PATH")

	// 设置默认值
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.username", "root")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.database", "article_analysis")
	viper.SetDefault("database.driver", "mysql")

	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")

	viper.SetDefault("openai.api_base", "https://api.moonshot.cn/v1")
	viper.SetDefault("openai.model", "kimi-k2-0905-preview")

	viper.SetDefault("log.level", "info")

	// 读取环境变量
	viper.AutomaticEnv()

	// 绑定特定的环境变量，设置优先级
	viper.BindEnv("openai.api_key", "OPENAI_API_KEY")
	viper.BindEnv("openai.api_base", "OPENAI_API_BASE")
	viper.BindEnv("openai.model", "OPENAI_MODEL")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件不存在，使用默认值
			fmt.Println("配置文件不存在，使用默认配置")
		} else {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	return &config, nil
}
