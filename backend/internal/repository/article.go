package repository

import (
	"article-analysis/internal/model"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

// ExistsByTitle 检查是否存在同标题文章
func (r *ArticleRepository) ExistsByTitle(title string) (bool, error) {
	var count int64
	if err := r.db.Model(&model.Article{}).Where("title = ?", title).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
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
		// 逃逸LIKE通配符，防止构造恶意模式
		escaped := escapeLike(req.Keyword)
		pattern := "%" + escaped + "%"
		query = query.Where("(title LIKE ? ESCAPE '\\' OR content LIKE ? ESCAPE '\\' OR author LIKE ? ESCAPE '\\')",
			pattern, pattern, pattern)
	}

	if req.Author != "" {
		query = query.Where("author = ?", req.Author)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 排序 - 使用白名单并通过Clause构造避免字符串拼接
	validSortColumns := map[string]bool{"title": true, "author": true, "upload_time": true}
	desc := strings.EqualFold(req.Order, "desc")
	if validSortColumns[req.Sort] {
		query = query.Clauses(clause.OrderBy{Columns: []clause.OrderByColumn{{
			Column: clause.Column{Name: req.Sort},
			Desc:   desc,
		}}})
	} else {
		query = query.Clauses(clause.OrderBy{Columns: []clause.OrderByColumn{{
			Column: clause.Column{Name: "upload_time"},
			Desc:   true,
		}}})
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	err := query.Offset(offset).Limit(req.PageSize).Find(&articles).Error
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
		escaped := escapeLike(req.Keyword)
		pattern := "%" + escaped + "%"
		query = query.Where("(a.title LIKE ? ESCAPE '\\' OR a.content LIKE ? ESCAPE '\\' OR a.author LIKE ? ESCAPE '\\')",
			pattern, pattern, pattern)
	}

	if req.Author != "" {
		query = query.Where("a.author = ?", req.Author)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 排序 - 使用白名单并通过Clause构造避免字符串拼接
	validSortColumns := map[string]bool{"title": true, "author": true, "upload_time": true}
	desc := strings.EqualFold(req.Order, "desc")
	if validSortColumns[req.Sort] {
		query = query.Clauses(clause.OrderBy{Columns: []clause.OrderByColumn{{
			Column: clause.Column{Name: "a." + req.Sort},
			Desc:   desc,
		}}})
	} else {
		query = query.Clauses(clause.OrderBy{Columns: []clause.OrderByColumn{{
			Column: clause.Column{Name: "a.upload_time"},
			Desc:   true,
		}}})
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	err := query.Offset(offset).Limit(req.PageSize).Scan(&articles).Error
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
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string    `gorm:"type:varchar(500);not null" json:"title"`
	Author      string    `gorm:"type:varchar(200);not null;index" json:"author"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	FilePath    string    `gorm:"type:varchar(500);not null" json:"file_path"`
	FileSize    int64     `gorm:"not null" json:"file_size"`
	WordCount   int       `gorm:"default:0" json:"word_count"`
	UploadTime  time.Time `gorm:"autoCreateTime" json:"upload_time"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	HasAnalysis bool      `gorm:"-" json:"has_analysis"` // 临时字段，不存储到数据库
}

// escapeLike 逃逸LIKE通配符，配合 ESCAPE '\\' 使用
func escapeLike(s string) string {
	// 先逃逸反斜杠，再处理通配符
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "%", "\\%")
	s = strings.ReplaceAll(s, "_", "\\_")
	return s
}

// ArticleWithAnalysis 包含分析状态的文章信息
type ArticleWithAnalysis struct {
	ID             uint64    `json:"id"`
	Title          string    `json:"title"`
	Author         string    `json:"author"`
	FilePath       string    `json:"file_path"`
	FileSize       int64     `json:"file_size"`
	UploadTime     time.Time `json:"upload_time"`
	CreatedAt      time.Time `json:"created_at"`
	AnalysisStatus string    `json:"analysis_status"`
	HasAnalysis    bool      `json:"has_analysis"`
}
