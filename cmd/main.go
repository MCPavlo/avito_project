package main

import (
	"fmt"
	"log"
	"net/http"

	"tz/internal/config"
	"tz/internal/db"
	"tz/internal/handler"
	"tz/internal/services"
)

func main() {
	cfg, err := config.LoadConfig("./internal/config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := db.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	userService := services.NewUserService(db)
	authHandler := handler.NewAuthHandler(userService)

	pvzService := services.NewPVZService(db)
	pvzHandler := handler.NewPVZHandler(pvzService)

	receptionService := services.NewReceptionService(db)
	receptionHandler := handler.NewReceptionHandler(receptionService)

	goodsService := services.NewGoodsService(db)
	goodsHandler := handler.NewGoodsHandler(goodsService)

	dummyLoginHandler := handler.NewDummyLoginHandler()

	http.HandleFunc("/register", authHandler.Register)
	http.HandleFunc("/login", authHandler.Login)
	http.HandleFunc("/pvz/create", pvzHandler.Create)
	http.HandleFunc("/receptions/create", receptionHandler.Create)
	http.HandleFunc("/pvz/{pvzId}/close_last_reception", receptionHandler.Close)
	http.HandleFunc("/products", goodsHandler.Add)
	http.HandleFunc("/pvz/{pvzId}/delete_last_product", goodsHandler.DeleteLast)
	http.HandleFunc("/pvz", pvzHandler.GetPVZs)
	http.HandleFunc("/receptions", receptionHandler.GetReceptions)
	http.HandleFunc("/dummyLogin", dummyLoginHandler.Login)

	fmt.Printf("Server is starting on port %s...\n", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(cfg.Server.Port, nil))
}
