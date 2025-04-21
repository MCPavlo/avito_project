package services

import (
	"errors"
	"tz/internal/db"
)

// GoodsService предоставляет методы для работы с товарами
type GoodsService struct {
	db *db.DB
}

// NewGoodsService создает новый экземпляр GoodsService
func NewGoodsService(db *db.DB) *GoodsService {
	return &GoodsService{db: db}
}

// Add добавляет товар в текущую приемку
func (gs *GoodsService) Add(pvzID int, productType string) (db.Goods, error) {
	// Получение последней открытой приемки
	var reception db.Reception
	err := gs.db.QueryRow("SELECT id, pvz_id, date_time, status FROM receptions WHERE pvz_id = $1 AND status = 'in_progress' ORDER BY date_time DESC LIMIT 1", pvzID).Scan(&reception.ID, &reception.PVZID, &reception.CreatedAt, &reception.Status)
	if err != nil {
		return db.Goods{}, errors.New("no open reception found")
	}

	// Добавление товара в приемку
	var product db.Goods
	err = gs.db.QueryRow("INSERT INTO products (type, reception_id, date_time) VALUES ($1, $2, NOW()) RETURNING id, type, reception_id, date_time", productType, reception.ID).Scan(&product.ID, &product.Type, &product.Reception, &product.ReceivedAt)
	if err != nil {
		return db.Goods{}, err
	}

	return product, nil
}

// DeleteLast удаляет последний добавленный товар из текущей приемки
func (gs *GoodsService) DeleteLast(pvzID int) (db.Goods, error) {
	// Получение последней открытой приемки
	var reception db.Reception
	err := gs.db.QueryRow("SELECT id, pvz_id, date_time, status FROM receptions WHERE pvz_id = $1 AND status = 'in_progress' ORDER BY date_time DESC LIMIT 1", pvzID).Scan(&reception.ID, &reception.PVZID, &reception.CreatedAt, &reception.Status)
	if err != nil {
		return db.Goods{}, errors.New("no open reception found")
	}

	// Получение последнего добавленного товара
	var product db.Goods
	err = gs.db.QueryRow("SELECT id, type, reception_id, date_time FROM products WHERE reception_id = $1 ORDER BY date_time DESC LIMIT 1", reception.ID).Scan(&product.ID, &product.Type, &product.ID, &product.ReceivedAt)
	if err != nil {
		return db.Goods{}, errors.New("no products to delete")
	}

	// Удаление товара
	_, err = gs.db.Exec("DELETE FROM products WHERE id = $1", product.ID)
	if err != nil {
		return db.Goods{}, err
	}

	return product, nil
}
