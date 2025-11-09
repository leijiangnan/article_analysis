package repository

import (
	"article-analysis/internal/model"
	"time"

	"gorm.io/gorm"
)

type ArticleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

func (r *ArticleRepository) Create(article *model.Article) error {
	return r.db.Create(article).Error
}

func (r *ArticleRepository) GetByID(id uint64) (*model.Article, error) {
	var article model.Article
	err := r.db.First(&article, id).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *ArticleRepository) GetList(req *model.PaginationRequest) (*model.PaginationResponse, error) {
	var articles []model.Article
	var total int64
	
	query := r.db.Model(&model.Article{})
	
	// 搜索条件
	if req.Keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ? OR author LIKE ?", 
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	
	if req.Author != "" {
		query = query.Where("author = ?", req.Author)
	}
	
	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	
	// 排序
	order := req.Sort + " " + req.Order
	if req.Sort != "title" && req.Sort != "author" && req.Sort != "upload_time" {
		order = "upload_time DESC"
	}
	
	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	err := query.Order(order).Offset(offset).Limit(req.PageSize).Find(&articles).Error
	if err != nil {
		return nil, err
	}
	
	return &model.PaginationResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		List:     articles,
	}, nil
}

// GetListWithAnalysis 获取文章列表及分析状态
func (r *ArticleRepository) GetListWithAnalysis(req *model.PaginationRequest) (*model.PaginationResponse, error) {
	var articles []ArticleWithAnalysis
	var total int64
	
	// 使用JOIN查询获取文章及其分析状态
	query := r.db.Table("articles a").
		Select(`a.id, a.title, a.author, a.file_path, a.file_size, 
			a.upload_time, a.created_at, 
			IFNULL(aa.analysis_status, 'none') as analysis_status,
			CASE WHEN aa.id IS NOT NULL THEN true ELSE false END as has_analysis`).
		Joins("LEFT JOIN article_analyses aa ON a.id = aa.article_id")
	
	// 搜索条件
	if req.Keyword != "" {
		query = query.Where("a.title LIKE ? OR a.content LIKE ? OR a.author LIKE ?", 
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	
	if req.Author != "" {
		query = query.Where("a.author = ?", req.Author)
	}
	
	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	
	// 排序
	order := req.Sort + " " + req.Order
	if req.Sort != "title" && req.Sort != "author" && req.Sort != "upload_time" {
		order = "a.upload_time DESC"
	}
	
	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	err := query.Order(order).Offset(offset).Limit(req.PageSize).Scan(&articles).Error
	if err != nil {
		return nil, err
	}
	
	return &model.PaginationResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		List:     articles,
	}, nil
}

func (r *ArticleRepository) Update(article *model.Article) error {
	return r.db.Save(article).Error
}

func (r *ArticleRepository) Delete(id uint64) error {
	return r.db.Delete(&model.Article{}, id).Error
}

func (r *ArticleRepository) GetAuthors() ([]string, error) {
	var authors []string
	err := r.db.Model(&model.Article{}).Distinct().Pluck("author", &authors).Error
	return authors, err
}

// 扩展Article结构体
func (a *Article) TableName() string {
	return "articles"
}

type Article struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Title      string    `gorm:"type:varchar(500);not null" json:"title"`
	Author     string    `gorm:"type:varchar(200);not null;index" json:"author"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	FilePath   string    `gorm:"type:varchar(500);not null" json:"file_path"`
	FileSize   int64     `gorm:"not null" json:"file_size"`
	WordCount  int       `gorm:"default:0" json:"word_count"`
	UploadTime time.Time `gorm:"autoCreateTime" json:"upload_time"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	HasAnalysis bool     `gorm:"-" json:"has_analysis"` // 临时字段，不存储到数据库
}

// ArticleWithAnalysis 包含分析状态的文章信息
type ArticleWithAnalysis struct {
	ID              uint64    `json:"id"`
	Title           string    `json:"title"`
	Author          string    `json:"author"`
	FilePath        string    `json:"file_path"`
	FileSize        int64     `json:"file_size"`
	UploadTime      time.Time `json:"upload_time"`
	CreatedAt       time.Time `json:"created_at"`
	AnalysisStatus  string    `json:"analysis_status"`
	HasAnalysis     bool      `json:"has_analysis"`
}