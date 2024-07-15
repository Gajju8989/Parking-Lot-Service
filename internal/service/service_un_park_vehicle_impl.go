package service

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"parking_lot_service/internal/genericresponse"
	"parking_lot_service/internal/repo/models"
	"parking_lot_service/internal/service/model"
	"time"
)

func (s *impl) UnParkVehicle(ctx context.Context, req *model.UnParkVehicleRequest) (
	*model.UnParkVehicleResponse, error) {

	parkedVehicle, err := s.parkingLotRepo.GetParkedVehicle(ctx, req.VehicleNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &genericresponse.GenericResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Record Not Found",
			}
		}
		return nil, &genericresponse.GenericResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	err = s.parkingLotRepo.DeleteParkedVehicle(ctx, parkedVehicle)
	if err != nil {
		return nil, &genericresponse.GenericResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Unable To Delete Parked Vehicle",
		}
	}

	// Get the count of available parking spots for the given ParkingLotID and VehicleID
	cnt, err := s.parkingLotRepo.GetAvailableParkingSpotsByParkingLotIdAndVehicleId(ctx,
		int(req.ParkingLotID), int(req.VehicleID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &genericresponse.GenericResponse{
				StatusCode: http.StatusNotFound,
				Message:    "Record Not Found",
			}
		}
		return nil, &genericresponse.GenericResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	// Get the Maximum Count By Request ParkingLotId  and  VehicleId.
	maxSpots, err := getMaxSpotsInParkingLot(req)
	if err != nil {
		return nil, err
	}

	// Check if the available spots exceed the maximum limit
	if cnt >= maxSpots {
		return nil, &genericresponse.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "All Spots Are Already Free for this Vehicle Type",
		}
	}

	// Update the parking space to increment the available spots
	updateParkingPayload := &models.ParkingSpace{
		ParkingLotId:   req.ParkingLotID,
		VehicleTypeId:  req.VehicleID,
		AvailableSpots: cnt + 1,
	}
	err = s.parkingLotRepo.UpdateParkingSpace(ctx, updateParkingPayload)
	if err != nil {
		// If there is an error updating the parking space, return an internal server error
		return nil, &genericresponse.GenericResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Unable to update parking space",
		}
	}

	// Calculate the fare and duration
	entryTime := parkedVehicle.EntryTime
	exitTime := time.Now()
	duration := exitTime.Sub(entryTime)

	totalFare, err := calculateFare(int(req.ParkingLotID), int(req.VehicleID), duration)

	if err != nil {
		return nil, &genericresponse.GenericResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	// Create the response with parking receipt details
	response := &model.UnParkVehicleResponse{
		Parking: model.ParkingReceipt{
			VehicleNumber: req.VehicleNumber,
			TotalFare:     totalFare,
			From:          entryTime.Format(time.RFC3339),
			To:            exitTime.Format(time.RFC3339),
			VehicleID:     int(req.VehicleID),
			ParkingLotID:  int(req.ParkingLotID),
		},
	}
	return response, nil

}

func calculateFare(parkingLotID int, vehicleTypeId int, duration time.Duration) (float64, error) {
	var tariffModels = map[int]map[int]model.Tariff{
		1: {
			1: {HourlyRate: 5},                                                       // Motorcycles/scooters
			2: {HourlyRate: 20.5},                                                    // Cars/SUVs
			3: {HourlyRate: 50, DayRate: 500, MaxDurationForDayRate: 24 * time.Hour}, // Buses/Trucks
		},
		2: {
			1: {HourlyRate: 10.5},                          // Motorcycles/scooters
			2: {FirstHourRate: 50, AdditionalHourRate: 25}, // Cars/SUVs
			3: {HourlyRate: 100},                           // Buses/Trucks
		},
	}

	tariff, ok := tariffModels[parkingLotID][vehicleTypeId]
	if !ok {
		return 0, fmt.Errorf("no tariff found for parking lot %d and vehicle type %d", parkingLotID, vehicleTypeId)
	}

	// Calculate the number of hours rounded up
	hours := int(duration.Hours())
	remainingMinutes := int(duration.Minutes()) % 60
	if remainingMinutes > 0 {
		hours++
	}

	// Calculate the fare based on the tariff model
	switch {
	case tariff.DayRate > 0 && duration <= tariff.MaxDurationForDayRate:
		// Case 1: If a day rate exists and the duration is within the max duration for the day rate
		{
			return float64(hours) * tariff.HourlyRate, nil
		}
	case tariff.DayRate > 0 && duration > tariff.MaxDurationForDayRate:
		// Case 2: If a day rate exists and the duration exceeds the max duration for the day rate
		{
			days := int(duration / tariff.MaxDurationForDayRate)
			remainingHours := duration - time.Duration(days)*tariff.MaxDurationForDayRate

			// Calculate the additional hours beyond the max day rate duration
			additionalHours := int(remainingHours.Hours())
			if remainingHours.Minutes() > float64(additionalHours*60) {
				additionalHours++
			}
			days--
			fare := (float64(days) * tariff.DayRate) + (float64(additionalHours) * tariff.HourlyRate) +
				(float64(tariff.MaxDurationForDayRate.Hours()) * tariff.HourlyRate)

			return fare, nil
		}
	case tariff.FirstHourRate > 0 && tariff.AdditionalHourRate > 0:
		// Case 3: If a special rate exists for the first hour and a different rate for additional hours
		{
			if hours == 1 {
				return tariff.FirstHourRate, nil
			}
			additionalHours := hours - 1
			return tariff.FirstHourRate + float64(additionalHours)*tariff.AdditionalHourRate, nil
		}
	default:
		// Case 4: Default case with a standard hourly rate
		{
			return float64(hours) * tariff.HourlyRate, nil
		}
	}
}

func getMaxSpotsInParkingLot(req *model.UnParkVehicleRequest) (int, error) {
	var maxSpots int

	switch {
	case req.ParkingLotID == 1 && req.VehicleID == 1:
		maxSpots = 50
	case req.ParkingLotID == 1 && req.VehicleID == 2:
		maxSpots = 30
	case req.ParkingLotID == 1 && req.VehicleID == 3:
		maxSpots = 20
	case req.ParkingLotID == 2 && req.VehicleID == 1:
		maxSpots = 100
	case req.ParkingLotID == 2 && req.VehicleID == 2:
		maxSpots = 80
	case req.ParkingLotID == 2 && req.VehicleID == 3:
		maxSpots = 40
	default:
		// If the ParkingLotID and VehicleID combination is invalid, return a bad request error
		return -1, &genericresponse.GenericResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Wrong Parking LotId or VehicleId",
		}
	}
	return maxSpots, nil
}
