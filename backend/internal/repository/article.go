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