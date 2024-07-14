package handler

import (
	"github.com/labstack/echo/v4"
	"parking_lot_service/internal/service"
)

type ParkingLotHandler interface {
	GetFreeParkingSpaces(c echo.Context) error
	GetParkingSpaceByParkingLotId(c echo.Context) error
	ParkVehicle(c echo.Context) error
}

type impl struct {
	parkingLotSvc service.ParkingLotService
}

func NewParkingLotHandler(parkingLotSvc service.ParkingLotService) ParkingLotHandler {
	return &impl{
		parkingLotSvc: parkingLotSvc,
	}
}
