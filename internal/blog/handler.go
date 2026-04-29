package blog

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func HandleGetBlogs(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	limitNum := 0

	if limit != "" {
		var err error
		limitNum, err = strconv.Atoi(limit)
		if err != nil || limitNum < 0 {
			http.Error(w, "Limit parametresi sorunlu", http.StatusBadRequest)
			return
		}
	}

	blogs, err := GetBlogs(limitNum)
	if err != nil {
		http.Error(w, "Veriler çekilemedi.", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blogs)
}

func HandleGetBlogById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	idNum, err := strconv.Atoi(id)
	if err != nil || idNum < 1 {
		http.Error(w, "Yanlış ID formatı..", http.StatusBadRequest)
		return
	}

	blog, err := GetBlogsById(idNum)
	if err != nil {
		http.Error(w, "Böyle bir blog yok.", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blog)
}

func HandleCreateBlogs(w http.ResponseWriter, r *http.Request) {
	var b Blog
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "JSON hatası.", http.StatusBadRequest)
		return
	}

	if err := CreateBlog(&b); err != nil {
		http.Error(w, "Kaydedilemedi.", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "başarılı!"})
}
