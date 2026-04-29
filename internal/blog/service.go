package blog

import (
	"demo/internal/db"

	"gorm.io/gorm"
)

func GetBlogs(page int, limit int, field []string) ([]Blog, error) {
	var blogs []Blog

	tx := db.DB.Session(&gorm.Session{})

	if page > 0 {
		pageSize := 10 // i wanted to define this variable myself.
		offset := (page - 1) * pageSize
		tx = tx.Offset(offset).Limit(pageSize)
	}

	if limit > 0 {
		tx = tx.Limit(limit)
	}

	if len(field) > 0 && field[0] != "" {
		tx = tx.Select(field)
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
