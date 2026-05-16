package blog

import (
	"demo/internal/middleware"
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
			http.Error(w, "Invalid page parameter.", http.StatusBadRequest)
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
			http.Error(w, "Invalid limit parameter.", http.StatusBadRequest)
			return
		}
		if limitNum > 100 {
			limitNum = 100
			log.Println(limitNum)
		}
	}

	// search
	searchQuery := r.URL.Query().Get("search")

	// field
	field := r.URL.Query().Get("field")

	if _, err := strconv.Atoi(field); err == nil {
		http.Error(w, "Invalid field parameter.", http.StatusBadRequest)
		return
	}

	arrOfField := strings.Split(field, ",")

	// sort
	sortQuery := r.URL.Query().Get("sort")
	lowerCaseSortQ := strings.ToLower(sortQuery)
	if _, err := strconv.Atoi(lowerCaseSortQ); err == nil || (lowerCaseSortQ != "" && lowerCaseSortQ != "desc" && lowerCaseSortQ != "asc") {
		http.Error(w, "Invalid sort parameter. Use 'asc' or 'desc'.", http.StatusBadRequest)
		return
	}

	// Response
	blogs, totalCount, filteredCount, err := GetBlogs(pageNum, limitNum, searchQuery, arrOfField, sortQuery)
	if err != nil {
		http.Error(w, "Failet to fetch blogs.", http.StatusInternalServerError)
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
		http.Error(w, "Invalid blog ID format.", http.StatusBadRequest)
		return
	}

	userIDVal := r.Context().Value(middleware.UserIDKey)
	var userID uint = 0

	if userIDVal != nil {
		userID = userIDVal.(uint)
	}

	blog, err := GetBlogsById(idNum, userID)
	if err != nil {
		http.Error(w, "Blog not found.", http.StatusNotFound)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blog)
}

func HandleCreateBlogs(w http.ResponseWriter, r *http.Request) {
	var blog Blog
	// the ID of who created the blog
	ctxID := r.Context().Value(middleware.UserIDKey)
	if ctxID == nil {
		http.Error(w, "Unauthorized: Token missing.", http.StatusUnauthorized)
		return
	}
	userID, ok := ctxID.(uint)
	if !ok {
		http.Error(w, "Internal context identity error.", http.StatusInternalServerError)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&blog); err != nil {
		http.Error(w, "Invalid JSON format.", http.StatusBadRequest)
		return
	}

	blog.AuthorID = userID

	if strings.TrimSpace(blog.Title) == "" || strings.TrimSpace(blog.Content) == "" {
		http.Error(w, "Title and content fields are required.", http.StatusBadRequest)
		return
	}

	if err := CreateBlog(&blog); err != nil {
		http.Error(w, "Failed to create blog post.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Blog post created successfully."})
}

func HandleUpdateBlog(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid blog ID format.", http.StatusBadRequest)
		return
	}

	var input Blog
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON format.", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(input.Title) == "" && strings.TrimSpace(input.Content) == "" && input.IsPublic == nil {
		http.Error(w, "Update body cannot be empty.", http.StatusBadRequest)
		return
	}

	UpdatedBlogDatas := Blog{
		Title:    strings.TrimSpace(input.Title),
		Content:  strings.TrimSpace(input.Content),
		IsPublic: input.IsPublic,
	}

	// the id of the user who tried to update the blog
	ctxID := r.Context().Value(middleware.UserIDKey)
	if ctxID == nil {
		http.Error(w, "Unauthorized: Token missing.", http.StatusUnauthorized)
		return
	}
	userID, ok := ctxID.(uint)
	if !ok {
		http.Error(w, "Internal context identity error.", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := UpdateBlog(&UpdatedBlogDatas, uint(id), uint(userID))
	if err != nil {
		http.Error(w, "Failed to update blog post.", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Blog not found or action forbidden.", http.StatusForbidden)
		return
	}

	log.Println(rowsAffected, "rows affected.")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Blog post updated successfully."})
}

func HandleDeleteBlog(w http.ResponseWriter, r *http.Request) {
	idQuery := r.PathValue("id")
	id, err := strconv.ParseUint(idQuery, 10, 32)
	if err != nil {
		http.Error(w, "Invalid blog ID format.", http.StatusBadRequest)
		return
	}

	// the id of the user who tried to delete the blog
	ctxID := r.Context().Value(middleware.UserIDKey)
	if ctxID == nil {
		http.Error(w, "Unauthorized: Token missing.", http.StatusUnauthorized)
		return
	}
	userID, ok := ctxID.(uint)
	if !ok {
		http.Error(w, "Internal context identity error.", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := DeleteBlog(uint(id), uint(userID))
	if err != nil {
		http.Error(w, "Failed to delete blog post.", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Blog not found or action forbidden.", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Blog post deleted successfully."})
}
