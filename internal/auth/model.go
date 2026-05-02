package auth

type User struct {
	ID       uint   `gorm:"primarykey" json:"id,omitempty"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Password string `gorm:"not null" json:"-"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}
