package router

import (
	"parking_lot_service/internal/handler"

	"github.com/labstack/echo/v4"
)

type Router interface {
	MapRoutes(e *echo.Echo)
}

type impl struct {
	parkingLotHandler handler.ParkingLotHandler
}

func NewRouter(parkingLotHandler handler.ParkingLotHandler) Router {
	return &impl{
		parkingLotHandler: parkingLotHandler,
	}
}
