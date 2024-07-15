package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"parking_lot_service/internal/genericresponse"
	"strconv"
)

// @Summary Get free parking spaces
// @Description Retrieve the number of free parking spaces in all parking lots
// @ID get-free-parking-spaces
// @Produce json
// @Failure 404,500 {object} genericresponse.GenericResponse
// @Success 200 {array} model.FreeSpotsResponse
// @Router /parking-lot/free-parking-spaces [get]
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

// @Summary Get parking space by parking lot ID
// @Description Retrieve the details of a specific parking space by its parking lot ID
// @ID get-parking-space-by-id
// @Param parking_lot_id query integer true "Parking Lot ID"
// @Produce json
// @Success 200 {object} model.FreeSpotsResponse
// @Failure 404,500 {object} genericresponse.GenericResponse
// @Router /parking-lot/parking-space [get]
func (s *impl) GetParkingSpaceByParkingLotId(c echo.Context) error {
	var (
		ctx                    = c.Request().Context()
		parkingLotIdQueryParam = c.QueryParam("parking_lot_id")
		parkingLotId, err      = strconv.Atoi(parkingLotIdQueryParam)
	)

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
