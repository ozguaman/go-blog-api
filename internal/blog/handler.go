package blog

import (
	"encoding/json"
	"log"
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
		if limitNum > 100 {
			limitNum = 100
			log.Println(limitNum)
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
	blogs, totalCount, filteredCount, err := GetBlogs(pageNum, limitNum, searchQuery, arrOfField, sortQuery)
	if err != nil {
		http.Error(w, "Veriler çekilemedi.", http.StatusInternalServerError)
		log.Println("GORM Hatası:", err)
		return
	}

	response := BlogResponse{
		TotalCount:    totalCount,
		FilteredCount: filteredCount,
		Response:      blogs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
	var blog Blog
	// the ID of who created the blog
	var userID = r.Context().Value("userID").(uint)

	if err := json.NewDecoder(r.Body).Decode(&blog); err != nil {
		http.Error(w, "JSON hatası.", http.StatusBadRequest)
		return
	}

	blog.AuthorID = userID

	if strings.TrimSpace(blog.Title) == "" || strings.TrimSpace(blog.Content) == "" {
		http.Error(w, "Title ve Content alanları boş olamaz.", http.StatusBadRequest)
		return
	}

	if err := CreateBlog(&blog); err != nil {
		http.Error(w, "Kaydedilemedi.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "başarılı!"})
}

func HandleUpdateBlog(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Yanlış formatta ID girişi.", http.StatusBadRequest)
		return
	}

	var input Blog
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Json hatası.", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(input.Title) == "" && strings.TrimSpace(input.Content) == "" {
		http.Error(w, "Title veya Contentten en az biri dolu olmalı.", http.StatusBadRequest)
		return
	}

	UpdatedBlogDatas := Blog{
		Title:   strings.TrimSpace(input.Title),
		Content: strings.TrimSpace(input.Content),
	}

	// the id of the user who tried to update the blog
	userID := r.Context().Value("userID").(uint)

	rowsAffected, err := UpdateBlog(&UpdatedBlogDatas, uint(id), uint(userID))
	if err != nil {
		http.Error(w, "Veri güncellenemedi.", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Böyle bir blog bulunmamakta veya yetkinizi aşıyor.", http.StatusBadRequest)
		return
	}

	log.Println(rowsAffected, " satır etkilendi.")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Güncelleme başarıyla tamamlandı."})
}

func HandleDeleteBlog(w http.ResponseWriter, r *http.Request) {
	idQuery := r.PathValue("id")
	id, err := strconv.ParseUint(idQuery, 10, 32)
	if err != nil {
		http.Error(w, "Yanlış id formatı.", http.StatusBadRequest)
		return
	}

	// the id of the user who tried to delete the blog
	userID := r.Context().Value("userID").(uint)

	rowsAffected, err := DeleteBlog(uint(id), uint(userID))
	if err != nil {
		http.Error(w, "Veri silinirken bir hata oluştu.", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Böyle bir blog yok veya yetkinizi aşıyor.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Veri başarıyla silindi."})
}
