package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "userID"

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

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userID := uint(claims["user_id"].(float64))

			ctx := context.WithValue(r.Context(), UserIDKey, userID)

			next(w, r.WithContext(ctx))
			return
		} else {
			http.Error(w, "Token bilgileri okunamadı.", http.StatusUnauthorized)
			return
		}
	}
}

func AuthOptionalMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			next(w, r)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			next(w, r)
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Beklenmeyen token formatı")
			}
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		if err == nil && token.Valid {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				userID := uint(claims["user_id"].(float64))
				ctx := context.WithValue(r.Context(), UserIDKey, userID)
				next(w, r.WithContext(ctx))
				return
			}
		}

		next(w, r)
	}
}
