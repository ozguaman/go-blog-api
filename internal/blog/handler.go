package blog

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func HandleGetBlogs(w http.ResponseWriter, r *http.Request) {

	// pagination
	page := r.URL.Query().Get("page")
	pageNum := 1

	if page != "" {
		var err error
		pageNum, err = strconv.Atoi(page)
		if err != nil || pageNum <= 0 {
			http.Error(w, "Page parametresi sorunlu.", http.StatusBadRequest)
			return
		}
	}

	// limit
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

	// search
	searchQuery := r.URL.Query().Get("search")
	if _, err := strconv.Atoi(searchQuery); err == nil {
		http.Error(w, "Bozuk formatta search parametresi girişi.", http.StatusBadRequest)
		return
	}

	// field
	field := r.URL.Query().Get("field")

	if _, err := strconv.Atoi(field); err == nil {
		http.Error(w, "Bozuk formatta field girişi.", http.StatusBadRequest)
		return
	}
	arrOfField := strings.Split(field, ",")

	// sort
	sortQuery := r.URL.Query().Get("sort")
	lowerCaseSortQ := strings.ToLower(sortQuery)
	if _, err := strconv.Atoi(lowerCaseSortQ); err == nil || (lowerCaseSortQ != "" && lowerCaseSortQ != "desc" && lowerCaseSortQ != "asc") {
		http.Error(w, "Yanlış veya eksik sort parametresi girişi.", http.StatusBadRequest)
		return
	}

	// Response
	blogs, err := GetBlogs(pageNum, limitNum, searchQuery, arrOfField, sortQuery)
	if err != nil {
		http.Error(w, "Veriler çekilemedi.", http.StatusInternalServerError)
		return
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

	if b.Title == "" || b.Content == "" {
		http.Error(w, "Title ve Content alanları boş olamaz.", http.StatusBadRequest)
		return
	}

	if err := CreateBlog(&b); err != nil {
		http.Error(w, "Kaydedilemedi.", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "başarılı!"})
}
