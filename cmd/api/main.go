package main

import (
	"demo/internal/blog"
	"demo/internal/db"
	"log"
	"net/http"
)

func main() {

	db.Connect()

	err := db.DB.AutoMigrate(&blog.Blog{})
	if err != nil {
		log.Fatal("Tablo oluşturulurken bir sorun oluştu. ", err)
	}

	http.HandleFunc("GET /blogs", blog.HandleGetBlogs)
	http.HandleFunc("GET /blogs/{id}", blog.HandleGetBlogById)
	http.HandleFunc("POST /blogs", blog.HandleCreateBlogs)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
