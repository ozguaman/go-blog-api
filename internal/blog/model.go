package blog

import "gorm.io/gorm"

type Blog struct {
	gorm.Model
	Title   string `json: "title"`
	Content string `json: "content"`
}
