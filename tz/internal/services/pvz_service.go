package services

import (
	"errors"
	"tz/internal/db"
)

// PVZService предоставляет методы для работы с ПВЗ
type PVZService struct {
	db *db.DB
}

// NewPVZService создает новый экземпляр PVZService
func NewPVZService(db *db.DB) *PVZService {
	return &PVZService{db: db}
}

// Create создает новый ПВЗ
func (ps *PVZService) Create(city string) (db.PVZ, error) {
	// Проверка допустимости города
	if city != "Москва" && city != "Санкт-Петербург" && city != "Казань" {
		return db.PVZ{}, errors.New("invalid city")
	}

	// Создание нового ПВЗ
	var pvz db.PVZ
	err := ps.db.QueryRow("INSERT INTO pvz (city, registration) VALUES ($1, NOW()) RETURNING id, city, registration", city).Scan(&pvz.ID, &pvz.City, &pvz.Registration)
	if err != nil {
		return db.PVZ{}, err
	}

	return pvz, nil
}

func (ps *PVZService) GetPVZs(startDate, endDate string, page, limit int) ([]db.PVZ, error) {

	// Проверка допустимости параметров
	if page < 1 || limit < 1 || limit > 30 {
		return nil, errors.New("invalid pagination parameters")
	}

	// Подготовка SQL-запроса
	query := "SELECT id, city, registration FROM pvz WHERE registration BETWEEN $1 AND $2 ORDER BY registration DESC LIMIT $3 OFFSET $4"

	// Выполнение запроса
	rows, err := ps.db.Query(query, startDate, endDate, limit, (page-1)*limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Сканирование результатов
	var pvzs []db.PVZ
	for rows.Next() {
		var pvz db.PVZ
		err := rows.Scan(&pvz.ID, &pvz.City, &pvz.Registration)
		if err != nil {
			return nil, err
		}
		pvzs = append(pvzs, pvz)
	}

	return pvzs, nil
}
