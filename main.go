package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"parking_lot_service/internal/database/postgresql/config"
	"parking_lot_service/internal/database/postgresql/migration"
	handler2 "parking_lot_service/internal/handler"
	"parking_lot_service/internal/repo"
	router2 "parking_lot_service/internal/router"
	"parking_lot_service/internal/service"
)

func main() {
	config.InitDB()
	err := migration.MigrateAll(config.GetDB())
	if err != nil {
		return
	}
	e := echo.New()
	db := repo.NewParkingLotRepo(config.GetDB())
	err = db.SeedParkingSpace(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	srvc := service.NewParkingLotService(db)
	handler := handler2.NewParkingLotHandler(srvc)
	router := router2.NewRouter(handler)
	fmt.Println("success")
	router.MapRoutes(e)
	port := ":8080"
	fmt.Printf("Server started on port %s\n", port)
	e.Logger.Fatal(e.Start(port))
}
