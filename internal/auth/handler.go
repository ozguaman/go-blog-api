package auth

import (
	"encoding/json"
	"net/http"
	"strings"
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

	user := User{
		Email:    registerRequest.Email,
		Username: registerRequest.Username,
		Password: registerRequest.Password,
	}

	if err := Register(&user); err != nil {
		http.Error(w, "Kullanıcı kaydı aşamasında bir sorun oluştu.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Kayıt işlemi başarılı."})
}
