package model

import (
	"time"
)

type Article struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Title      string    `gorm:"type:varchar(500);not null" json:"title"`
	Author     string    `gorm:"type:varchar(200);not null;index" json:"author"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	FilePath   string    `gorm:"type:varchar(500);not null" json:"file_path"`
	FileSize   int64     `gorm:"not null" json:"file_size"`
	UploadTime time.Time `gorm:"autoCreateTime" json:"upload_time"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ArticleAnalysis struct {
	ID               uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ArticleID        uint64    `gorm:"not null;index" json:"article_id"`
	CoreViewpoints   string    `gorm:"type:text" json:"core_viewpoints"`
	FileStructure    string    `gorm:"type:text" json:"file_structure"`
	AuthorThoughts   string    `gorm:"type:text" json:"author_thoughts"`
	RelatedMaterials string    `gorm:"type:text" json:"related_materials"`
	AnalysisStatus   string    `gorm:"type:enum('pending','processing','completed','failed');default:'pending'" json:"analysis_status"`
	AnalysisTime     *time.Time `json:"analysis_time"`
	ErrorMessage     string    `gorm:"type:text" json:"error_message"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	
	Article Article `gorm:"foreignKey:ArticleID" json:"article,omitempty"`
}

type PaginationRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=10" binding:"min=1,max=100"`
	Keyword  string `form:"keyword"`
	Author   string `form:"author"`
	Sort     string `form:"sort,default=upload_time"`
	Order    string `form:"order,default=desc" binding:"oneof=asc desc"`
}

type PaginationResponse struct {
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	List     interface{} `json:"list"`
}

type ApiResponse struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp"`
}