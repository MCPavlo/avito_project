package services

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"tz/internal/db"
)

type PVZDetail struct {
	PVZ        db.PVZ            `json:"pvz"`
	Receptions []ReceptionDetail `json:"receptions"`
}

type ReceptionDetail struct {
	Reception db.Reception `json:"reception"`
	Goods     []db.Goods   `json:"goods"`
}

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

func (ps *PVZService) GetPVZs(startDateStr, endDateStr string, page, limit int) ([]PVZDetail, error) {
	// Проверка корректности параметров пагинации
	if page < 1 || limit < 1 || limit > 30 {
		return nil, errors.New("invalid pagination parameters")
	}

	// Определяем, применять ли фильтр по дате, и парсим даты если они заданы
	var filterByDate bool
	var startTime, endTime time.Time
	if startDateStr != "" && endDateStr != "" {
		var err error
		startTime, err = time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid startDate format, expected RFC3339: %w", err)
		}
		endTime, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid endDate format, expected RFC3339: %w", err)
		}
		filterByDate = true
	}

	offset := (page - 1) * limit

	var pvzRows *sql.Rows
	var err error

	// Если активен фильтр по дате, выбираем ПВЗ, у которых есть хотя бы одна приёмка в заданном диапазоне
	if filterByDate {
		query := `
            SELECT DISTINCT p.id, p.city, p.registration
            FROM pvz p
            JOIN receptions r ON r.pvz_id = p.id
            WHERE r.created_at BETWEEN $1 AND $2
            ORDER BY p.registration DESC
            LIMIT $3 OFFSET $4
        `
		pvzRows, err = ps.db.Query(query, startTime, endTime, limit, offset)
	} else {
		query := `
            SELECT id, city, registration
            FROM pvz
            ORDER BY registration DESC
            LIMIT $1 OFFSET $2
        `
		pvzRows, err = ps.db.Query(query, limit, offset)
	}
	if err != nil {
		return nil, err
	}
	defer pvzRows.Close()

	var details []PVZDetail
	for pvzRows.Next() {
		var p db.PVZ
		if err := pvzRows.Scan(&p.ID, &p.City, &p.Registration); err != nil {
			return nil, err
		}

		// Для каждого ПВЗ выбираем связанные приёмки
		var recRows *sql.Rows
		if filterByDate {
			recQuery := `
                SELECT id, pvz_id, created_at, status
                FROM receptions
                WHERE pvz_id = $1 AND created_at BETWEEN $2 AND $3
                ORDER BY created_at DESC
            `
			recRows, err = ps.db.Query(recQuery, p.ID, startTime, endTime)
		} else {
			recQuery := `
                SELECT id, pvz_id, created_at, status
                FROM receptions
                WHERE pvz_id = $1
                ORDER BY created_at DESC
            `
			recRows, err = ps.db.Query(recQuery, p.ID)
		}
		if err != nil {
			return nil, err
		}

		var receptions []ReceptionDetail
		for recRows.Next() {
			var rec db.Reception
			if err := recRows.Scan(&rec.ID, &rec.PVZID, &rec.CreatedAt, &rec.Status); err != nil {
				recRows.Close()
				return nil, err
			}

			// Для каждой приёмки выбираем связанные товары
			goodsQuery := `
                SELECT id, type, received_at
                FROM goods
                WHERE reception_id = $1
                ORDER BY id
            `
			goodsRows, err := ps.db.Query(goodsQuery, rec.ID)
			if err != nil {
				recRows.Close()
				return nil, err
			}
			var goodsList []db.Goods
			for goodsRows.Next() {
				var g db.Goods
				if err := goodsRows.Scan(&g.ID, &g.Type, &g.ReceivedAt); err != nil {
					goodsRows.Close()
					recRows.Close()
					return nil, err
				}
				// Привязываем родительскую приёмку к товару
				g.Reception = rec
				goodsList = append(goodsList, g)
			}
			goodsRows.Close()

			receptions = append(receptions, ReceptionDetail{
				Reception: rec,
				Goods:     goodsList,
			})
		}
		recRows.Close()

		details = append(details, PVZDetail{
			PVZ:        p,
			Receptions: receptions,
		})
	}
	if err = pvzRows.Err(); err != nil {
		return nil, err
	}

	return details, nil
}
