package services

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/dgrijalva/jwt-go"
	"log"
	"time"
	_ "time"
	"tz/internal/db"
)

// UserService предоставляет методы для работы с пользователями
type UserService struct {
	db *db.DB
}

// NewUserService создает новый экземпляр UserService
func NewUserService(db *db.DB) *UserService {
	return &UserService{db: db}
}

// Register регистрирует нового пользователя
func (us *UserService) Register(email, password, role string) error {
	// Проверка существования пользователя с таким email
	var existingUser db.User
	err := us.db.Get(&existingUser, "SELECT * FROM users WHERE email = $1", email)
	if err == nil {
		return errors.New("user with this email already exists")
	}

	// Создание нового пользователя
	_, err = us.db.Exec("INSERT INTO users (email, password, role) VALUES ($1, $2, $3)", email, password, role)
	if err != nil {
		return err
	}

	return nil
}

// Login аутентифицирует пользователя
func (us *UserService) Login(email, password string) (string, error) {
	// Проверка существования пользователя с таким email и паролем
	var user db.User
	err := us.db.QueryRow("SELECT id, email, password, role FROM users WHERE email = $1 AND password = $2", email, password).Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		log.Printf("Ошибка поиска пользователя: %v", err)
		return "", errors.New("invalid credentials")
	}

	// Генерация JWT-токена
	token, err := generateJWT(user.ID, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func generateJWT(userID int, role string) (string, error) {
	// Создаем JWT-токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"role":   role,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // Токен истекает через 24 часа
	})

	// Подписываем токен секретным ключом
	secretKey := []byte("secret_key")
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
