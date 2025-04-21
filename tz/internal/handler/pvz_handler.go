package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"tz/internal/services"
)

// PVZHandler предоставляет методы для работы с ПВЗ
type PVZHandler struct {
	pvzService *services.PVZService
}

// NewPVZHandler создает новый экземпляр PVZHandler
func NewPVZHandler(pvzService *services.PVZService) *PVZHandler {
	return &PVZHandler{pvzService: pvzService}
}

// Create создает новый ПВЗ
func (ph *PVZHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Проверка роли пользователя (например, через JWT-токен)
	userRole, err := getUserRoleFromJWT(r) // Реализуйте эту функцию для получения роли пользователя из JWT

	if userRole != "moderator" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	var pvz struct {
		City string `json:"city"`
	}

	err = json.NewDecoder(r.Body).Decode(&pvz)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdPVZ, err := ph.pvzService.Create(pvz.City)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdPVZ)
}

// GetPVZs получает список ПВЗ с фильтрацией по дате и пагинацией
func (ph *PVZHandler) GetPVZs(w http.ResponseWriter, r *http.Request) {

	userRole, err := getUserRoleFromJWT(r) // Реализуйте эту функцию для получения роли пользователя из JWT

	if userRole != "moderator" && userRole != "employee" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Получение параметров запроса
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	// Получение списка ПВЗ
	pvzs, err := ph.pvzService.GetPVZs(startDate, endDate, page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pvzs)
}
