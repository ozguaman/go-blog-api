package blog

import (
	"demo/internal/db"
)

func GetBlogs(userID uint, page int, limit int, searchQ string, field []string, sortQ string) ([]Blog, int64, int64, error) {
	var blogs []Blog
	var totalCount, filteredCount int64

	db.DB.Count(&totalCount) // getting count of the all blogs

	if page > 0 {
		pageSize := 10
		offset := (page - 1) * pageSize
		db.DB = db.DB.Offset(offset).Limit(pageSize)
	}

	if limit > 0 {
		db.DB = db.DB.Limit(limit)
	}

	if searchQ != "" {
		query := "%" + searchQ + "%"
		db.DB = db.DB.Where("title ILIKE ? OR content ILIKE ?", query, query)
	}

	if len(field) > 0 && field[0] != "" {
		db.DB = db.DB.Select(field)
	}

	if sortQ != "" {
		db.DB = db.DB.Order("created_at " + sortQ)
	}

	result := db.DB.Where("is_public = ? or author_id = ?", true, userID).Find(&blogs)
	db.DB.Count(&filteredCount) // getting count of the filtered blogs
	return blogs, totalCount, filteredCount, result.Error
}

func GetBlogsById(id int, userID uint) (Blog, error) {
	var blog Blog

	result := db.DB.Where("is_public = ? or author_id = ?", true, userID).First(&blog, id)
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
