package service

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"net/http"
	"parking_lot_service/internal/genericresponse"
	"parking_lot_service/internal/service/model"
)

func (s *impl) GetFreeParkingSpaces(ctx context.Context) ([]*model.FreeSpotsResponse, error) {

	resp, err := s.parkingLotRepo.GetParkingSpaces(ctx)
	if err != nil || len(resp) == 0 {
		if errors.Is(err, gorm.ErrRecordNotFound) || len(resp) == 0 {
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

	var (
		motorcycleSpots = 0
		carSuvSpots     = 0
		busTruckSpots   = 0
	)

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
