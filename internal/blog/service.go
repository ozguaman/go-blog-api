package blog

import (
	"demo/internal/db"

	"gorm.io/gorm"
)

func GetBlogs(userID uint, page int, limit int, searchQ string, field []string, sortQ string) ([]Blog, int64, int64, error) {
	var blogs []Blog
	var totalCount, filteredCount int64

	tx := db.DB.Session(&gorm.Session{}).Model(&Blog{})

	tx.Count(&totalCount) // getting count of the all blogs

	if page > 0 {
		pageSize := 10
		offset := (page - 1) * pageSize
		tx = tx.Offset(offset).Limit(pageSize)
	}

	if limit > 0 {
		tx = tx.Limit(limit)
	}

	if searchQ != "" {
		query := "%" + searchQ + "%"
		tx = tx.Where("title ILIKE ? OR content ILIKE ?", query, query)
	}

	if len(field) > 0 && field[0] != "" {
		tx = tx.Select(field)
	}

	if sortQ != "" {
		tx = tx.Order("created_at " + sortQ)
	}

	result := tx.Where("is_public = ? or author_id = ?", true, userID).Find(&blogs)
	tx.Count(&filteredCount) // getting count of the filtered blogs
	return blogs, totalCount, filteredCount, result.Error
}

func GetBlogsById(id int, userID uint) (Blog, error) {
	var blog Blog

	tx := db.DB.Session(&gorm.Session{})

	result := tx.Where("is_public = ? or author_id = ?", true, userID).First(&blog, id)
	return blog, result.Error
}

func CreateBlog(blog *Blog) error {
	return db.DB.Create(blog).Error
}

func UpdateBlog(blog *Blog, id uint, userID uint) (int64, error) {
	result := db.DB.Model(&Blog{}).Where("id = ? and author_id = ?", id, userID).Updates(blog)
	return result.RowsAffected, result.Error
}

func DeleteBlog(id uint, userID uint) (int64, error) {
	var blog *Blog
	// db.DB.Unscoped().Where("id = ?", id).Delete(&Blog{}) -> it does disable soft delete.
	result := db.DB.Where("id = ? and author_id = ?", id, userID).Delete(&blog)
	return result.RowsAffected, result.Error
}
