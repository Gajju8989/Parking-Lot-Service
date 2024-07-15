package router

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "parking_lot_service/docs"
)

func (r *impl) MapRoutes(e *echo.Echo) {
	parkingLot := e.Group("/parking-lot")
	parkingLot.GET("/free-parking-spaces", r.parkingLotHandler.GetFreeParkingSpaces)
	parkingLot.GET("/parking-space", r.parkingLotHandler.GetParkingSpaceByParkingLotId)
	parkingLot.POST("/park-vehicle", r.parkingLotHandler.ParkVehicle)
	parkingLot.POST("/un-park-vehicle", r.parkingLotHandler.UnParkVehicle)

	// Swagger endpoint
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
