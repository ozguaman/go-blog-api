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

func GetBlogsById(id int) (Blog, error) {
	var blog Blog

	tx := db.DB.Session(&gorm.Session{})

	result := tx.First(&blog, id)
	return blog, result.Error
}

func CreateBlog(b *Blog) error {
	return db.DB.Create(b).Error
}
