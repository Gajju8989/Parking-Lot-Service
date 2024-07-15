package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"parking_lot_service/internal/genericresponse"
	"parking_lot_service/internal/service/model"
)

// @Summary Park a vehicle
// @Description Park a vehicle in the parking lot based on the provided details
// @ID park-vehicle
// @Accept json
// @Produce json
// @Param request body model.ParkVehicleRequest true "Vehicle details to park"
// @Success 200 {object} model.ParkVehicleResponse
// @Failure 400,404,500 {object} genericresponse.GenericResponse
// @Router /parking-lot/park-vehicle [post]
func (s *impl) ParkVehicle(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		req = &model.ParkVehicleRequest{}
		err = c.Bind(&req)
	)

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
