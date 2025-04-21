package db

import "time"

// User представляет пользователя
type User struct {
	ID       int
	Email    string
	Password string
	Role     string
}

// PVZ представляет пункт выдачи заказов
type PVZ struct {
	ID           int
	City         string
	Registration time.Time
}

// Reception представляет приемку товаров
type Reception struct {
	ID        int
	PVZID     int
	CreatedAt time.Time
	Status    string
}

// Goods представляет товар
type Goods struct {
	ID         int
	Type       string
	ReceivedAt time.Time
	Reception  Reception
}
