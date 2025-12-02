package main

import (
	"article-analysis/internal/config"
	"article-analysis/internal/handler"
	"article-analysis/internal/middleware"
	"article-analysis/internal/repository"
	"article-analysis/internal/service"
	"article-analysis/pkg/logger"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 加载环境文件（优先级：.env -> .env.local；同时尝试父目录）
	// 忽略不存在的文件加载错误
	_ = godotenv.Load(".env")
	_ = godotenv.Load(".env.local")
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("../.env.local")

	// 初始化配置
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	log := logger.NewLogger(cfg.Log.Level)
	defer log.Sync()

	// 输出OpenAI配置加载状态（不打印敏感信息）
	log.Info("OpenAI配置已加载",
		zap.String("api_base", cfg.OpenAI.APIBase),
		zap.String("model", cfg.OpenAI.Model),
		zap.Bool("has_api_key", cfg.OpenAI.APIKey != ""),
	)

	// 连接数据库
	db, err := initDB(cfg)
	if err != nil {
		log.Fatal("数据库连接失败", zap.Error(err))
		os.Exit(1)
	}

	// 自动迁移
	if err := autoMigrate(db); err != nil {
		log.Fatal("数据库迁移失败", zap.Error(err))
		os.Exit(1)
	}

	// 初始化依赖
	articleRepo := repository.NewArticleRepository(db)
	analysisRepo := repository.NewAnalysisRepository(db)

	articleService := service.NewArticleService(articleRepo, log)
	analysisService := service.NewAnalysisService(analysisRepo, articleRepo, cfg, log)

	articleHandler := handler.NewArticleHandler(articleService)
	analysisHandler := handler.NewAnalysisHandler(analysisService)

	// 设置路由
	router := setupRouter(articleHandler, analysisHandler, log)

	// 启动服务
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Info("服务器启动", zap.String("address", addr))
	if err := router.Run(addr); err != nil {
		log.Fatal("服务器启动失败", zap.Error(err))
		os.Exit(1)
	}
}

func initDB(cfg *config.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.Database.Driver {
	case "sqlite":
		dbPath := "./data/article_analysis.db"
		// 确保数据目录存在
		if err := os.MkdirAll("./data", 0755); err != nil {
			return nil, fmt.Errorf("创建数据目录失败: %w", err)
		}
		dialector = sqlite.Open(dbPath)
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Database.Username,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Database,
		)
		dialector = mysql.Open(dsn)
	default:
		return nil, fmt.Errorf("不支持的数据库驱动: %s", cfg.Database.Driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&repository.Article{},
		&repository.ArticleAnalysis{},
	)
}

func setupRouter(articleHandler *handler.ArticleHandler, analysisHandler *handler.AnalysisHandler, log *logger.Logger) *gin.Engine {
	router := gin.New()

	// 全局中间件
	router.Use(middleware.Logger(log))
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	// 健康检查端点
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "article-analysis-backend",
		})
	})

	// 静态文件服务
	router.Static("/uploads", "./web/uploads")

	// API路由组
	api := router.Group("/api")
	{
		// 文章相关路由
		articles := api.Group("/articles")
		{
			articles.POST("/upload", articleHandler.UploadArticle)
			articles.POST("/create", articleHandler.CreateArticle)
			articles.GET("/authors", articleHandler.GetAuthors)
			articles.GET("", articleHandler.GetArticleList)
			articles.GET("/with-analysis", articleHandler.GetArticleListWithAnalysis)
			articles.GET("/:id", articleHandler.GetArticleDetail)
			articles.DELETE("/:id", articleHandler.DeleteArticle)
			articles.POST("/:id/analyze", analysisHandler.AnalyzeArticle)
			articles.GET("/:id/analysis", analysisHandler.GetAnalysisResult)
		}

		// 分析任务状态
		api.GET("/analysis/status/:task_id", analysisHandler.GetAnalysisStatus)
	}

	return router
}
