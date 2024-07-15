package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"parking_lot_service/internal/genericresponse"
	"parking_lot_service/internal/service/model"
)

// @Summary Unpark a vehicle
// @Description Remove a parked vehicle from the parking lot based on the provided details
// @ID unpark-vehicle
// @Accept json
// @Produce json
// @Param request body model.UnParkVehicleRequest true "Vehicle details to unpark"
// @Success 200 {object} model.UnParkVehicleResponse
// @Failure 400,404,500 {object} genericresponse.GenericResponse
// @Router /parking-lot/un-park-vehicle [post]
func (s *impl) UnParkVehicle(c echo.Context) error {
	var (
		ctx = c.Request().Context()
		req = &model.UnParkVehicleRequest{}
		err = c.Bind(&req)
	)

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
