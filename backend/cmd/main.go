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

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 初始化配置
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	log := logger.NewLogger(cfg.Log.Level)
	defer log.Sync()

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

	// 静态文件服务
	router.Static("/uploads", "./web/uploads")

	// API路由组
	api := router.Group("/api/v1")
	{
		// 文章相关路由
		articles := api.Group("/articles")
		{
			articles.POST("/upload", articleHandler.UploadArticle)
			articles.GET("/authors", articleHandler.GetAuthors)
			articles.GET("", articleHandler.GetArticleList)
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