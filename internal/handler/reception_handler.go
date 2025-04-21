package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"tz/internal/services"
)

// ReceptionHandler предоставляет методы для работы с приемками
type ReceptionHandler struct {
	receptionService *services.ReceptionService
}

// NewReceptionHandler создает новый экземпляр ReceptionHandler
func NewReceptionHandler(receptionService *services.ReceptionService) *ReceptionHandler {
	return &ReceptionHandler{receptionService: receptionService}
}

// Create создает новую приемку
func (rh *ReceptionHandler) Create(w http.ResponseWriter, r *http.Request) {

	userRole, err := getUserRoleFromJWT(r)

	if userRole != "employee" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	var reception struct {
		PVZID int `json:"pvzId"`
	}

	err = json.NewDecoder(r.Body).Decode(&reception)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdReception, err := rh.receptionService.Create(reception.PVZID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdReception)
}

// Close закрывает последнюю открытую приемку
func (rh *ReceptionHandler) Close(w http.ResponseWriter, r *http.Request) {

	userRole, err := getUserRoleFromJWT(r)

	if userRole != "employee" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	var reception struct {
		PVZID int `json:"pvzId"`
	}

	err = json.NewDecoder(r.Body).Decode(&reception)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	closedReception, err := rh.receptionService.Close(reception.PVZID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(closedReception)
}

// GetReceptions получает список приемок с фильтрацией по дате и пагинацией
func (rh *ReceptionHandler) GetReceptions(w http.ResponseWriter, r *http.Request) {

	userRole, err := getUserRoleFromJWT(r)

	if userRole != "employee" && userRole != "moderator" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}
	
	// Получение параметров запроса
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	// Получение списка приемок
	receptions, err := rh.receptionService.GetReceptions(startDate, endDate, page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(receptions)
}
