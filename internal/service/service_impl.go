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

func (s *impl) GetFreeParkingSpaces(ctx context.Context) ([]*model.FreeSpotsResponse, error) {

	resp, err := s.parkingLotRepo.GetParkingSpaces(ctx)
	if err != nil {
		return nil, handleCommonErrors(err)
	}
	var freeSpotsResponses []*model.FreeSpotsResponse

	for i := 1; i <= 2; i++ {
		motorcycleSpots := 0
		carSuvSpots := 0
		busTruckSpots := 0
		for _, parkingSpace := range resp {
			if i == int(parkingSpace.ParkingLotId) {
				switch parkingSpace.VehicleTypeId {
				case 1: // Motorcycles/scooters
					motorcycleSpots += parkingSpace.AvailableSpots
				case 2: // Cars/SUVs
					carSuvSpots += parkingSpace.AvailableSpots
				case 3: // Buses/Trucks
					busTruckSpots += parkingSpace.AvailableSpots
				}
			}
		}
		freeSpotsResponses = append(freeSpotsResponses, &model.FreeSpotsResponse{
			ParkingLotID:                    i,
			FreeSpotsForMotorcyclesScooters: motorcycleSpots,
			FreeSpotsForCarsSUVs:            carSuvSpots,
			FreeSpotsForBusesTrucks:         busTruckSpots,
		})
	}
	return freeSpotsResponses, nil
}

func (s *impl) GetFreeParkingSpaceById(ctx context.Context, parkingLotId int) (*model.FreeSpotsResponse, error) {

	resp, err := s.parkingLotRepo.GetFreeParkingSpaceById(ctx, parkingLotId)

	if err != nil || len(resp) == 0 {
		if errors.Is(err, gorm.ErrRecordNotFound) || len(resp) == 0 {
			return nil, &genericresponse.GenericResponse{
				StatusCode: http.StatusNotFound,
				Message:    "parking lot not found",
			}
		}
		return nil, &genericresponse.GenericResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	motorcycleSpots := 0
	carSuvSpots := 0
	busTruckSpots := 0
	for _, parkingSpace := range resp {
		switch parkingSpace.VehicleTypeId {
		case 1: // Motorcycles/scooters
			motorcycleSpots += parkingSpace.AvailableSpots
		case 2: // Cars/SUVs
			carSuvSpots += parkingSpace.AvailableSpots
		case 3: // Buses/Trucks
			busTruckSpots += parkingSpace.AvailableSpots
		}
	}

	return &model.FreeSpotsResponse{
		ParkingLotID:                    parkingLotId,
		FreeSpotsForMotorcyclesScooters: motorcycleSpots,
		FreeSpotsForCarsSUVs:            carSuvSpots,
		FreeSpotsForBusesTrucks:         busTruckSpots,
	}, nil
}

func (s *impl) ParkVehicle(ctx context.Context, req *model.ParkVehicleRequest) (*model.ParkVehicleResponse, error) {

	cnt, err := s.parkingLotRepo.GetAvailableParkingSpotsByParkingLotIdAndVehicleId(ctx,
		int(req.ParkingLotID), int(req.VehicleID))

	if err != nil {
		return nil, handleCommonErrors(err)
	}

	if cnt > 0 {
		updateParkingPayload := &models.ParkingSpace{
			ParkingLotId:   req.ParkingLotID,
			VehicleTypeId:  req.VehicleID,
			AvailableSpots: cnt - 1,
		}
		err = s.parkingLotRepo.UpdateParkingSpace(ctx, updateParkingPayload)
		if err != nil {
			return nil, &genericresponse.GenericResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Unable to update parking space",
			}
		}
	} else {
		return nil, &genericresponse.GenericResponse{
			StatusCode: http.StatusNotFound,
			Message:    "No Spots Available ",
		}
	}

	err = s.parkingLotRepo.SaveParkedVehicle(ctx, &models.ParkedVehicle{
		VehicleNumber: req.VehicleNumber,
		ParkingLotID:  req.ParkingLotID,
		VehicleTypeId: req.VehicleID,
		VehicleName:   req.VehicleName,
		EntryTime:     time.Now(),
	})

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, &genericresponse.GenericResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Vehicle already in parking space",
			}
		}
		return nil, &genericresponse.GenericResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	resp := &model.ParkVehicleResponse{
		ParkingTicket: model.ParkingTicket{
			VehicleNumber: req.VehicleNumber,
			VehicleID:     int(req.VehicleID),
			EntryTime:     time.Now(),
		}}

	switch req.ParkingLotID {
	case 1:
		{
			resp.ParkingTicket.ParkingLot = "Parking Lot A"
		}
	case 2:
		{
			resp.ParkingTicket.ParkingLot = "Parking Lot B"
		}
	}

	return resp, nil
}

func (s *impl) UnParkVehicle(ctx context.Context, req *model.UnParkVehicleRequest) (
	*model.UnParkVehicleResponse, error) {

	parkedVehicle, err := s.parkingLotRepo.GetParkedVehicle(ctx, req.VehicleNumber)
	if err != nil {
		return nil, handleCommonErrors(err)
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
		return nil, handleCommonErrors(err)
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
	// CalculateFare calculates the parking fare based on the tariff model, vehicle type, and duration.

	tariff, ok := tariffModels[parkingLotID][vehicleTypeId]
	if !ok {
		return 0, fmt.Errorf("no tariff found for parking lot %d and vehicle type %d", parkingLotID, vehicleTypeId)
	}

	// Calculate the number of hours rounded up
	hours := int(duration.Hours())
	if duration.Minutes() > 0 {
		hours++
	}

	// Calculate the fare based on the tariff model
	switch {
	case tariff.DayRate > 0 && duration <= tariff.MaxDurationForDayRate:
		return float64(hours) * tariff.HourlyRate, nil
	case tariff.DayRate > 0 && duration > tariff.MaxDurationForDayRate:
		days := int(duration / tariff.MaxDurationForDayRate)
		remainingHours := duration - time.Duration(days)*tariff.MaxDurationForDayRate
		remainingFullHours := int(remainingHours.Hours())
		if remainingHours.Hours() > float64(remainingFullHours) {
			remainingFullHours++
		}
		return float64(days)*tariff.DayRate + float64(remainingFullHours)*tariff.HourlyRate, nil
	case tariff.FirstHourRate > 0 && tariff.AdditionalHourRate > 0:
		if hours == 1 {
			return tariff.FirstHourRate, nil
		}
		return tariff.FirstHourRate + float64(hours-1)*tariff.AdditionalHourRate, nil
	default:
		return float64(hours) * tariff.HourlyRate, nil
	}

}

func handleCommonErrors(err error) *genericresponse.GenericResponse {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &genericresponse.GenericResponse{
			StatusCode: http.StatusNotFound,
			Message:    err.Error(),
		}
	}
	return &genericresponse.GenericResponse{
		StatusCode: http.StatusInternalServerError,
		Message:    err.Error(),
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
