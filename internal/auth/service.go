package auth

import "demo/internal/db"

func Register(user *User) error {
	return db.DB.Create(user).Error
}
