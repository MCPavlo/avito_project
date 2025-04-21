package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

// DummyLoginHandler предоставляет метод для получения тестового токена
type DummyLoginHandler struct{}

// NewDummyLoginHandler создает новый экземпляр DummyLoginHandler
func NewDummyLoginHandler() *DummyLoginHandler {
	return &DummyLoginHandler{}
}

// Login возвращает тестовый токен с желаемым уровнем доступа
func (dlh *DummyLoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	var login struct {
		Role string `json:"role"`
	}

	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Генерация JWT-токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": login.Role,
	})

	// Подписываем токен секретным ключом
	secretKey := []byte("secret_key")
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
