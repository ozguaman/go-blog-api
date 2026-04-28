package blog

import (
	"encoding/json"
	"net/http"
)

func HandleGetBlogs(w http.ResponseWriter, r *http.Request) {
	blogs := GetBlogs()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blogs)
}
