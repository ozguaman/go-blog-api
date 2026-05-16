package blog

import (
	"demo/internal/db"
)

func GetBlogs(page int, limit int, searchQ string, field []string, sortQ string) ([]Blog, int64, int64, error) {
	var blogs []Blog
	var totalCount, filteredCount int64

	baseQuery := db.DB.Model(&Blog{})

	baseQuery.Count(&totalCount) // get the count of the all blogs

	if len(field) > 0 && field[0] != "" {
		baseQuery = baseQuery.Select(field)
	}

	if searchQ != "" {
		query := "%" + searchQ + "%"
		baseQuery = baseQuery.Where("title ILIKE ? OR content ILIKE ?", query, query)
	}

	baseQuery.Count(&filteredCount) // getting count of the filtered blogs

	if page > 0 {
		pageSize := 10
		offset := (page - 1) * pageSize
		baseQuery = baseQuery.Offset(offset).Limit(pageSize)
	}

	if limit > 0 {
		baseQuery = baseQuery.Limit(limit)
	}

	if sortQ != "" {
		baseQuery = baseQuery.Order("created_at " + sortQ)
	}

	result := baseQuery.Where("is_public = ?", true).Find(&blogs)
	return blogs, totalCount, filteredCount, result.Error
}

func GetBlogsById(id int, userID uint) (Blog, error) {
	var blog Blog

	baseQuery := db.DB.Model(&Blog{})

	result := baseQuery.Where("id = ? and (is_public = ? or author_id = ?)", id, true, userID).First(&blog)
	return blog, result.Error
}

func CreateBlog(blog *Blog) error {
	baseQuery := db.DB.Model(&Blog{})
	return baseQuery.Create(blog).Error
}

func UpdateBlog(blog *Blog, id uint, userID uint) (int64, error) {
	baseQuery := db.DB.Model(&Blog{})
	result := baseQuery.Model(&Blog{}).Where("id = ? and author_id = ?", id, userID).Updates(blog)
	return result.RowsAffected, result.Error
}

func DeleteBlog(id uint, userID uint) (int64, error) {
	var blog *Blog
	baseQuery := db.DB.Model(&Blog{})
	// baseQuery.Unscoped().Where("id = ?", id).Delete(&Blog{}) -> it does disable soft delete.
	result := baseQuery.Where("id = ? and author_id = ?", id, userID).Delete(&blog)
	return result.RowsAffected, result.Error
}
