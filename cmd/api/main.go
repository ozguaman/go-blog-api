package main

import (
	"demo/internal/auth"
	"demo/internal/blog"
	"demo/internal/db"
	"demo/internal/middleware"
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
	http.HandleFunc("GET /blogs", middleware.AuthMiddleware(blog.HandleGetBlogs))
	http.HandleFunc("GET /blogs/{id}", middleware.AuthMiddleware(blog.HandleGetBlogById))
	http.HandleFunc("POST /blogs", middleware.AuthMiddleware(blog.HandleCreateBlogs))
	http.HandleFunc("PATCH /blogs/{id}", middleware.AuthMiddleware(blog.HandleUpdateBlog))
	http.HandleFunc("DELETE /blogs/{id}", middleware.AuthMiddleware(blog.HandleDeleteBlog))

	// auth
	http.HandleFunc("POST /register", auth.HandleRegister)
	http.HandleFunc("POST /login", auth.HandleLogin)
	http.HandleFunc("PATCH /users/{id}", middleware.AuthMiddleware(auth.HandleUpdateUser))
	http.HandleFunc("DELETE /users/{id}", middleware.AuthMiddleware(auth.HandleDeleteUser))

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
