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
