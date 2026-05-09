package blog

import (
	"time"

	"gorm.io/gorm"
)

type Blog struct {
	ID        uint            `gorm:"primarykey" json:"id,omitempty"`
	AuthorID  uint            `json:"author_id,omitempty"`
	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Title     string          `json:"title,omitempty"`
	Content   string          `json:"content,omitempty"`
	IsPublic  *bool           `gorm:"default:true" json:"is_public,omitempty"`
}

type BlogResponse struct {
	TotalCount    int64  `json:"total_count"`
	FilteredCount int64  `json:"filtered_count"`
	Response      []Blog `json:"data"`
}
