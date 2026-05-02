package main

import (
	"demo/internal/auth"
	"demo/internal/blog"
	"demo/internal/db"
	"log"
	"net/http"
)

func main() {

	db.Connect()

	err := db.DB.AutoMigrate(&blog.Blog{}, &auth.User{})
	if err != nil {
		log.Fatal(err)
	}

	// blog
	http.HandleFunc("GET /blogs", blog.HandleGetBlogs)
	http.HandleFunc("GET /blogs/{id}", blog.HandleGetBlogById)
	http.HandleFunc("POST /blogs", blog.HandleCreateBlogs)
	http.HandleFunc("PATCH /blogs/{id}", blog.HandleUpdateBlog)
	http.HandleFunc("DELETE /blogs/{id}", blog.HandleDeleteBlog)

	// auth
	http.HandleFunc("POST /register", auth.HandleRegister)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
