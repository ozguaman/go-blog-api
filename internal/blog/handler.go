package blog

import (
	"encoding/json"
	"log"
	"net/http"
)

func HandleGetBlogs(w http.ResponseWriter, r *http.Request) {
	blogs, err := GetBlogs()
	if err != nil {
		log.Fatal("Dosyalar fetch edilirken bir sorun oldu. ", err)
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
