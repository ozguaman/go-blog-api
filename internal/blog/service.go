package blog

import (
	"demo/internal/db"
)

func GetBlogs() ([]Blog, error) {
	var blogs []Blog
	result := db.DB.Find(&blogs)
	return blogs, result.Error
}

func CreateBlog(b *Blog) error {
	return db.DB.Create(b).Error
}
