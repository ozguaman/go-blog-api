package auth

import (
	"demo/internal/blog"
	"demo/internal/db"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func Register(user *User) error {
	return db.DB.Create(user).Error
}

func FindUserByUsername(loginRequest LoginRequest) (User, error) {
	var user User
	result := db.DB.Where("username = ?", strings.TrimSpace(loginRequest.Username)).First(&user)
	return user, result.Error
}

func CreateToken(userId uint) (string, error) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env dosyası yüklenemedi!")
	}

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

	result := db.DB.Where("(author_id = ? AND is_public = ?) OR (author_id = ? and author_id = ?)", userID, true, requestID, userID).Find(&userBlogs)
	db.DB.Count(&totalCount)
	return userBlogs, totalCount, result.Error
}

func UpdateUser(idParam uint, user *User) (int64, error) {
	result := db.DB.Model(&User{}).Where("id = ?", idParam).Updates(&user)
	return result.RowsAffected, result.Error
}

func DeleteUser(userID uint64) (int64, error) {
	var rowsAffected int64

	err := db.DB.Transaction(func(tx *gorm.DB) error {

		if err := tx.Where("author_id = ?", userID).Delete(&blog.Blog{}); err.Error != nil {
			return err.Error
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
