package handlers

import (
	"encoding/json"
	"net/http"
	"note_service/internal/authorization"
	"note_service/internal/storage"

	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	Storage *storage.PostgreSQL
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Функция выдает токен для пользователя
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// 1️⃣ Ищем пользователя в БД
	user, err := h.Storage.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// 2️⃣ Проверяем пароль
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// 3️⃣ Генерируем токен с РЕАЛЬНЫМ ID
	token, err := authorization.GenerateToken(user.Id)
	if err != nil {
		http.Error(w, "could not generate token", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"token": token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
