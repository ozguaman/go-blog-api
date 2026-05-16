package auth

import (
	"demo/internal/blog"
	"demo/internal/db"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func Register(user *User) error {
	baseQuery := db.DB.Model(&User{})
	return baseQuery.Create(user).Error
}

func FindUserByUsername(loginRequest LoginRequest) (User, error) {
	var user User
	baseQuery := db.DB.Model(&User{})
	result := baseQuery.Where("username = ?", strings.TrimSpace(loginRequest.Username)).First(&user)
	return user, result.Error
}

func CreateToken(userId uint) (string, error) {

	ENV_JWT := os.Getenv("JWT_SECRET_KEY")
	JWT_SECRET_KEY := []byte(ENV_JWT)

	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(JWT_SECRET_KEY)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserByUserID(userID uint, requestID uint, page int, limit int, searchQ string, field []string, sortQ string) ([]blog.Blog, int64, error) {
	var userBlogs []blog.Blog
	var totalCount int64

	baseQuery := db.DB.Model(&blog.Blog{})

	baseQuery.Where("(author_id = ? AND is_public = ?) OR (author_id = ? AND author_id = ?)", userID, true, requestID, userID).Count(&totalCount)

	if searchQ != "" {
		query := "%" + searchQ + "%"
		baseQuery = baseQuery.Where("title ILIKE ? OR content ILIKE ?", query, query)
	}

	if page > 0 {
		pageSize := 10
		offset := (page - 1) * pageSize
		baseQuery = baseQuery.Offset(offset).Limit(pageSize)
	}

	if limit > 0 {
		baseQuery = baseQuery.Limit(limit)
	}

	if len(field) > 0 && field[0] != "" {
		baseQuery = baseQuery.Select(field)
	}

	if sortQ != "" {
		baseQuery = baseQuery.Order("created_at " + sortQ)
	}

	result := baseQuery.Where("(author_id = ? AND is_public = ?) OR (author_id = ?)", userID, true, requestID).Find(&userBlogs)
	return userBlogs, totalCount, result.Error
}

func UpdateUser(idParam uint, user *User) (int64, error) {
	baseQuery := db.DB.Model(&User{})
	result := baseQuery.Model(&User{}).Where("id = ?", idParam).Updates(&user)
	return result.RowsAffected, result.Error
}

func DeleteUser(userID uint64) (int64, error) {
	var rowsAffected int64

	err := db.DB.Transaction(func(tx *gorm.DB) error {

		blogResult := tx.Where("author_id = ?", uint(userID)).Delete(&blog.Blog{})
		if blogResult.Error != nil {
			return blogResult.Error
		}

		result := tx.Where("id = ?", userID).Delete(&User{})
		if result.Error != nil {
			return result.Error
		}

		rowsAffected = result.RowsAffected
		return nil
	})
	return rowsAffected, err
}
