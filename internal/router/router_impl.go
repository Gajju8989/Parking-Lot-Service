package router

import "github.com/labstack/echo/v4"

func (r *impl) MapRoutes(e *echo.Echo) {
	parkingLot := e.Group("/parking-lot")
	parkingLot.GET("/free-parking-spaces", r.parkingLotHandler.GetFreeParkingSpaces)
	parkingLot.GET("/free-parking-spaces", r.parkingLotHandler.GetParkingSpaceByParkingLotId)
	parkingLot.POST("/park-vehicle", r.parkingLotHandler.ParkVehicle)
	parkingLot.POST("/un-park-vehicle", r.parkingLotHandler.UnParkVehicle)
}
