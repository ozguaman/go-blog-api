package main

import (
	"demo/internal/blog"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("GET /blogs", blog.HandleGetBlogs)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
