package auth

import (
	"encoding/json"
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

func HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	var userUpdateRequest User
	requestID := r.Context().Value("userID").(uint)
	idParam := r.PathValue("id")

	idNum, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "Yanlış ID girişi.", http.StatusBadRequest)
		return
	}

	if idNum != uint64(requestID) {
		http.Error(w, "Böyle bir yetkiniz yok.", http.StatusForbidden)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&userUpdateRequest); err != nil {
		http.Error(w, "JSON hatası.", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(strings.TrimSpace(userUpdateRequest.Password)), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Sunucu kaynaklı bir güvenlik sorunu oluştu.", http.StatusInternalServerError)
		return
	}

	user := User{
		Email:    strings.TrimSpace(userUpdateRequest.Email),
		Username: strings.TrimSpace(userUpdateRequest.Username),
		Password: string(hashedPassword),
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
		http.Error(w, "JSON hatası.", http.StatusBadRequest)
		return
	}

	requestID := r.Context().Value("userID").(uint)

	if requestID != uint(idNum) {
		http.Error(w, "Yetkisiz işlem.", http.StatusUnauthorized)
		return
	}

	rowsAffected, err := DeleteUser(idNum)
	if err != nil {
		http.Error(w, "Kullanıcı silinemedi.", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Böyle bir kullanıcı yok.", http.StatusForbidden)
		return
	}
}
