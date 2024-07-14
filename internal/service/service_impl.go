package service

import (
	"context"
	"errors"
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &genericresponse.GenericResponse{
				StatusCode: http.StatusNotFound,
				Message:    err.Error(),
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

	err := s.parkingLotRepo.DecreaseAvailableSpot(ctx, &models.ParkingSpace{
		ParkingLotId:  req.ParkingLotID,
		VehicleTypeId: req.VehicleID,
	})

	if err != nil {
		return nil, &genericresponse.GenericResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Unable to update parking space",
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

	return nil, nil
}
