package handler

import (
	"encoding/json"
	"net/http"

	"tz/internal/services"
)

// GoodsHandler предоставляет методы для работы с товарами
type GoodsHandler struct {
	goodsService *services.GoodsService
}

// NewGoodsHandler создает новый экземпляр GoodsHandler
func NewGoodsHandler(goodsService *services.GoodsService) *GoodsHandler {
	return &GoodsHandler{goodsService: goodsService}
}

// Add добавляет товар в текущую приемку
func (gh *GoodsHandler) Add(w http.ResponseWriter, r *http.Request) {

	userRole, err := getUserRoleFromJWT(r)

	if userRole != "employee" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	var product struct {
		Type  string `json:"type"`
		PVZID int    `json:"pvzId"`
	}

	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdProduct, err := gh.goodsService.Add(product.PVZID, product.Type)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdProduct)
}

// DeleteLast удаляет последний добавленный товар из текущей приемки
func (gh *GoodsHandler) DeleteLast(w http.ResponseWriter, r *http.Request) {

	userRole, err := getUserRoleFromJWT(r)

	if userRole != "employee" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	var product struct {
		PVZID int `json:"pvzId"`
	}

	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	deletedProduct, err := gh.goodsService.DeleteLast(product.PVZID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deletedProduct)
}
