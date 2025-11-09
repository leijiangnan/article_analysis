package repository

import (
	"article-analysis/internal/model"
	"time"

	"gorm.io/gorm"
)

type AnalysisRepository struct {
	db *gorm.DB
}

func NewAnalysisRepository(db *gorm.DB) *AnalysisRepository {
	return &AnalysisRepository{db: db}
}

func (r *AnalysisRepository) Create(analysis *model.ArticleAnalysis) error {
	return r.db.Create(analysis).Error
}

func (r *AnalysisRepository) GetByArticleID(articleID uint64) (*model.ArticleAnalysis, error) {
	var analysis model.ArticleAnalysis
	err := r.db.Preload("Article").Where("article_id = ?", articleID).First(&analysis).Error
	if err != nil {
		return nil, err
	}
	return &analysis, nil
}

func (r *AnalysisRepository) Update(analysis *model.ArticleAnalysis) error {
	return r.db.Save(analysis).Error
}

func (r *AnalysisRepository) UpdateStatus(articleID uint64, status string, errorMsg string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"analysis_status": status,
		"error_message":   errorMsg,
		"updated_at":      now,
	}
	
	if status == "completed" || status == "failed" {
		updates["analysis_time"] = &now
	}
	
	return r.db.Model(&model.ArticleAnalysis{}).
		Where("article_id = ?", articleID).
		Updates(updates).Error
}

func (r *AnalysisRepository) GetByID(id uint64) (*model.ArticleAnalysis, error) {
	var analysis model.ArticleAnalysis
	err := r.db.First(&analysis, id).Error
	if err != nil {
		return nil, err
	}
	return &analysis, nil
}

// 扩展ArticleAnalysis结构体
func (a *ArticleAnalysis) TableName() string {
	return "article_analyses"
}

type ArticleAnalysis struct {
	ID               uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ArticleID        uint64     `gorm:"not null;index" json:"article_id"`
	CoreViewpoints   string     `gorm:"type:text" json:"core_viewpoints"`
	FileStructure    string     `gorm:"type:text" json:"file_structure"`
	AuthorThoughts   string     `gorm:"type:text" json:"author_thoughts"`
	RelatedMaterials string     `gorm:"type:text" json:"related_materials"`
	AnalysisStatus   string     `gorm:"default:'pending'" json:"analysis_status"`
	AnalysisTime     *time.Time `json:"analysis_time"`
	ErrorMessage     string     `gorm:"type:text" json:"error_message"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}