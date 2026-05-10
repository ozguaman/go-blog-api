package auth

import (
	"demo/internal/db"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
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

func UpdateUser(idParam int, user *User) error {
	result := db.DB.Model(&User{}).Where("id = ?", idParam).Updates(&user)
	return result.Error
}
