package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

const jwtSecret = "secret_key"

func getUserRoleFromJWT(r *http.Request) (string, error) {
	// Получаем заголовок Authorization
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}

	// Проверяем формат "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid authorization header format")
	}
	tokenString := parts[1]

	// Парсим токен с проверкой подписи
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем алгоритм подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	// Обрабатываем ошибки парсинга
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return "", errors.New("token expired")
			}
		}
		return "", fmt.Errorf("token parsing failed: %v", err)
	}

	// Проверяем валидность токена
	if !token.Valid {
		return "", errors.New("invalid token")
	}

	// Извлекаем claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	// Получаем роль
	role, ok := claims["role"].(string)
	if !ok || role == "" {
		return "", errors.New("role claim is missing or invalid")
	}

	log.Printf("User role from token: %s", role)

	return role, nil
}
