package services

import (
	"errors"
	"tz/internal/db"
)

// ReceptionService предоставляет методы для работы с приемками
type ReceptionService struct {
	db *db.DB
}

// NewReceptionService создает новый экземпляр ReceptionService
func NewReceptionService(db *db.DB) *ReceptionService {
	return &ReceptionService{db: db}
}

// Create создает новую приемку
func (rs *ReceptionService) Create(pvzID int) (db.Reception, error) {
	// Проверка существования незакрытой приемки
	var existingReception db.Reception
	err := rs.db.QueryRow("SELECT id FROM receptions WHERE pvz_id = $1 AND status = 'in_progress'", pvzID).Scan(&existingReception.ID)
	if err == nil {
		return db.Reception{}, errors.New("there is already an open reception")
	}

	// Создание новой приемки
	var reception db.Reception
	err = rs.db.QueryRow("INSERT INTO receptions (pvz_id, created_at, status) VALUES ($1, NOW(), 'in_progress') RETURNING id, pvz_id, created_at, status", pvzID).Scan(&reception.ID, &reception.PVZID, &reception.CreatedAt, &reception.Status)
	if err != nil {
		return db.Reception{}, err
	}

	return reception, nil
}

// Close закрывает последнюю открытую приемку
func (rs *ReceptionService) Close(pvzID int) (db.Reception, error) {
	// Получение последней открытой приемки
	var reception db.Reception
	err := rs.db.QueryRow("SELECT id, pvz_id, created_at, status FROM receptions WHERE pvz_id = $1 AND status = 'in_progress' ORDER BY created_at DESC LIMIT 1", pvzID).Scan(&reception.ID, &reception.PVZID, &reception.CreatedAt, &reception.Status)
	if err != nil {
		return db.Reception{}, errors.New("no open reception found")
	}

	// Закрытие приемки
	_, err = rs.db.Exec("UPDATE receptions SET status = 'close' WHERE id = $1", reception.ID)
	if err != nil {
		return db.Reception{}, err
	}

	return reception, nil
}

// GetReceptions получает список приемок с фильтрацией по дате и пагинацией
func (rs *ReceptionService) GetReceptions(startDate, endDate string, page, limit int) ([]db.Reception, error) {
	// Проверка допустимости параметров
	if page < 1 || limit < 1 || limit > 30 {
		return nil, errors.New("invalid pagination parameters")
	}

	// Подготовка SQL-запроса
	query := "SELECT id, pvz_id, created_at, status FROM receptions WHERE created_at BETWEEN $1 AND $2 ORDER BY created_at DESC LIMIT $3 OFFSET $4"

	// Выполнение запроса
	rows, err := rs.db.Query(query, startDate, endDate, limit, (page-1)*limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Сканирование результатов
	var receptions []db.Reception
	for rows.Next() {
		var reception db.Reception
		err := rows.Scan(&reception.ID, &reception.PVZID, &reception.CreatedAt, &reception.Status)
		if err != nil {
			return nil, err
		}
		receptions = append(receptions, reception)
	}

	return receptions, nil
}
