package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		authHader := r.Header.Get("Authorization")
		if authHader == "" {
			http.Error(w, "Yetkisiz işlem.", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Yanlış bilet formatı.", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Beklenmeyen token formatı: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Geçersiz veya süresi dolmuş token.", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
