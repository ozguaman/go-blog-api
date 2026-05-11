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

	tx := db.DB.Session(&gorm.Session{}).Model(&blog.Blog{})

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

	result := tx.Where("(author_id = ? AND is_public = ?) OR (author_id = ? and author_id = ?)", userID, true, requestID, userID).Find(&userBlogs)
	tx.Count(&totalCount)
	return userBlogs, totalCount, result.Error
}

func UpdateUser(idParam uint, user *User) (int64, error) {
	result := db.DB.Model(&User{}).Where("id = ?", idParam).Updates(&user)
	return result.RowsAffected, result.Error
}

func DeleteUser(userID uint64) (int64, error) {
	var user *User
	// blogları da sildirmen lazım unutma.
	result := db.DB.Where("id = ?", userID).Delete(&user)
	return result.RowsAffected, result.Error
}
