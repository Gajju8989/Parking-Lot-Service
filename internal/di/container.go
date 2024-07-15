package di

import (
	"context"
	"github.com/labstack/echo/v4"
	"parking_lot_service/internal/database/postgresql/config"
	"parking_lot_service/internal/database/postgresql/migration"
	handler2 "parking_lot_service/internal/handler"
	"parking_lot_service/internal/repo"
	router2 "parking_lot_service/internal/router"
	"parking_lot_service/internal/service"
)

// Container struct holds references to all dependencies
type Container struct {
	echoInstance *echo.Echo
	db           repo.ParkingLotRepo
}

// NewContainer initializes and returns a new Container instance
func NewContainer() *Container {
	config.InitDB()
	err := migration.MigrateAll(config.GetDB())
	if err != nil {
		return nil
	}

	e := echo.New()
	db := repo.NewParkingLotRepo(config.GetDB())
	err = db.SeedParkingSpace(context.Background())
	if err != nil {
		return nil
	}
	return &Container{
		echoInstance: e,
		db:           db,
	}
}

func (c *Container) GetEchoInstance() *echo.Echo {
	return c.echoInstance
}

func (c *Container) GetParkingLotRepo() repo.ParkingLotRepo {
	return c.db
}

func (c *Container) GetHandler() handler2.ParkingLotHandler {
	srvc := service.NewParkingLotService(c.db)
	return handler2.NewParkingLotHandler(srvc)
}

func (c *Container) GetRouter() router2.Router {
	handler := c.GetHandler()
	return router2.NewRouter(handler)
}
