package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"parking_lot_service/internal/genericresponse"
	"parking_lot_service/internal/service/model"
	"strconv"
)

func (s *impl) GetFreeParkingSpaces(c echo.Context) error {

	ctx := c.Request().Context()

	resp, err := s.parkingLotSvc.GetFreeParkingSpaces(ctx)
	if err != nil {
		// Handle specific errors if there is any genericError
		genericErr, ok := err.(*genericresponse.GenericResponse)
		if ok {
			return c.JSON(genericErr.StatusCode, genericErr)
		}
		// For any other unexpected errors, return a generic internal server error.
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")

	}

	return c.JSON(http.StatusOK, resp)
}

func (s *impl) GetParkingSpaceByParkingLotId(c echo.Context) error {
	ctx := c.Request().Context()
	parkingLotIdQueryParam := c.QueryParam("parking_lot_id")
	parkingLotId, err := strconv.Atoi(parkingLotIdQueryParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Parking lot id should be a number")
	}
	resp, err := s.parkingLotSvc.GetFreeParkingSpaceById(ctx, parkingLotId)
	if err != nil {
		// Handle specific errors if there is any genericError
		genericErr, ok := err.(*genericresponse.GenericResponse)
		if ok {
			return c.JSON(genericErr.StatusCode, genericErr)
		}
		// For any other unexpected errors, return a generic internal server error.
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	return c.JSON(http.StatusOK, resp)
}

func (s *impl) ParkVehicle(c echo.Context) error {
	ctx := c.Request().Context()
	req := &model.ParkVehicleRequest{}
	err := c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	resp, err := s.parkingLotSvc.ParkVehicle(ctx, req)
	if err != nil {
		genericErr, ok := err.(*genericresponse.GenericResponse)
		if ok {
			return c.JSON(genericErr.StatusCode, genericErr)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, resp)
}

func (s *impl) UnParkVehicle(c echo.Context) error {
	ctx := c.Request().Context()
	req := &model.UnParkVehicleRequest{}
	err := c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	resp, err := s.parkingLotSvc.UnParkVehicle(ctx, req)
	if err != nil {
		genericErr, ok := err.(*genericresponse.GenericResponse)
		if ok {
			return c.JSON(genericErr.StatusCode, genericErr)

		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	return c.JSON(http.StatusOK, resp)
}
