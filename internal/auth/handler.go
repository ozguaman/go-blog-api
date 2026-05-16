package auth

import (
	"demo/internal/blog"
	"demo/internal/middleware"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	var registerRequest RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
		http.Error(w, "JSON hatası.", http.StatusInternalServerError)
		return
	}

	if strings.TrimSpace(registerRequest.Email) == "" || strings.TrimSpace(registerRequest.Username) == "" || strings.TrimSpace(registerRequest.Password) == "" {
		http.Error(w, "Hiçbir alan boş bırakılamaz.", http.StatusBadRequest)
		return
	}

	// hash the password

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(strings.TrimSpace(registerRequest.Password)), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Sunucu kaynaklı bir güvenlik sorunu oluştu.", http.StatusInternalServerError)
		return
	}

	user := User{
		Email:    strings.TrimSpace(registerRequest.Email),
		Username: strings.TrimSpace(registerRequest.Username),
		Password: string(hashedPassword),
	}

	if err := Register(&user); err != nil {
		http.Error(w, "Kullanıcı kaydı aşamasında bir sorun oluştu.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Kayıt işlemi başarılı."})
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var loginRequest LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "JSON hatası.", http.StatusBadRequest)
		return
	}

	user, err := FindUserByUsername(loginRequest)
	if err != nil {
		http.Error(w, "Kullanıcı adı veya şifre yanlış.", http.StatusUnauthorized)
		return
	}

	password := strings.TrimSpace(loginRequest.Password)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		http.Error(w, "Kullanıcı adı veya şifre yanlış.", http.StatusUnauthorized)
		return
	}

	token, err := CreateToken(user.ID)

	message := map[string]string{
		"message": "Giriş başarılı.",
		"token":   token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(message)
}

func HandleGetBlogsByUserID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	idNum, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		http.Error(w, "Yanlış ID formatı.", http.StatusBadRequest)
		return
	}

	var requestID uint
	if ctxID := r.Context().Value(middleware.UserIDKey); ctxID != nil {
		requestID = ctxID.(uint)
	}

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

	blogs, totalCount, err := GetUserByUserID(uint(idNum), uint(requestID), pageNum, limitNum, searchQuery, arrOfField, sortQuery)
	if err != nil {
		http.Error(w, "Veriler çekilemedi.", http.StatusInternalServerError)
		return
	}

	response := blog.BlogResponse{
		TotalCount: totalCount,
		Response:   blogs,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	var userUpdateRequest User
	idParam := r.PathValue("id")

	idNum, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "Yanlış ID girişi.", http.StatusBadRequest)
		return
	}

	requestID := r.Context().Value(middleware.UserIDKey).(uint)

	if idNum != uint64(requestID) {
		http.Error(w, "Böyle bir yetkiniz yok.", http.StatusForbidden)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&userUpdateRequest); err != nil {
		http.Error(w, "JSON hatası.", http.StatusBadRequest)
		return
	}

	user := User{
		Email:    strings.TrimSpace(userUpdateRequest.Email),
		Username: strings.TrimSpace(userUpdateRequest.Username),
	}

	if strings.TrimSpace(userUpdateRequest.Password) != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(strings.TrimSpace(userUpdateRequest.Password)), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Sunucu kaynaklı bir güvenlik sorunu oluştu.", http.StatusInternalServerError)
			return
		}
		user.Password = string(hashedPassword)
	}

	rowsAffected, err := UpdateUser(uint(idNum), &user)
	if err != nil {
		http.Error(w, "Kullanıcı bilgileri güncellenemedi.", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Böyle bir kullanıcı yok.", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "İşlem başarılı!"})
}

func HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	idNum, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		http.Error(w, "Yanlış ID formatı.", http.StatusBadRequest)
		return
	}

	requestID := r.Context().Value(middleware.UserIDKey).(uint)

	if requestID != uint(idNum) {
		http.Error(w, "Yetkisiz işlem.", http.StatusUnauthorized)
		return
	}

	rowsAffected, err := DeleteUser(uint64(idNum))
	if err != nil {
		http.Error(w, "Kullanıcı silinemedi.", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Böyle bir kullanıcı yok.", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "İşlem başarılı."})
}
