package blog

import (
	"demo/internal/db"

	"gorm.io/gorm"
)

func GetBlogs(limit int) ([]Blog, error) {
	var blogs []Blog

	tx := db.DB.Session(&gorm.Session{})

	if limit > 0 {
		tx = tx.Limit(limit)
	}

	result := tx.Find(&blogs)
	return blogs, result.Error
}

func CreateBlog(b *Blog) error {
	return db.DB.Create(b).Error
}
