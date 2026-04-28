package main

import (
	"demo/internal/blog"
	"demo/internal/db"
	"log"
	"net/http"
)

func main() {

	db.Connect()

	http.HandleFunc("GET /blogs", blog.HandleGetBlogs)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
