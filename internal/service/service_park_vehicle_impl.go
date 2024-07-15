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

func (s *impl) ParkVehicle(ctx context.Context, req *model.ParkVehicleRequest) (*model.ParkVehicleResponse, error) {
	// Check available parking spots for the given parking lot and vehicle type
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

	// If there are available spots, update the parking space
	if cnt > 0 {
		// Decrease available spots by 1
		updateParkingPayload := &models.ParkingSpace{
			ParkingLotId:   req.ParkingLotID,
			VehicleTypeId:  req.VehicleID,
			AvailableSpots: cnt - 1,
		}
		// Update the parking space in the repository
		err = s.parkingLotRepo.UpdateParkingSpace(ctx, updateParkingPayload)
		if err != nil {
			return nil, &genericresponse.GenericResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "Unable to update parking space",
			}
		}
	} else {
		// No available spots, return error response
		return nil, &genericresponse.GenericResponse{
			StatusCode: http.StatusNotFound,
			Message:    "No Spots Available",
		}
	}

	// Save the parked vehicle details
	err = s.parkingLotRepo.SaveParkedVehicle(ctx, &models.ParkedVehicle{
		VehicleNumber: req.VehicleNumber,
		ParkingLotID:  req.ParkingLotID,
		VehicleTypeId: req.VehicleID,
		VehicleName:   req.VehicleName,
		EntryTime:     time.Now(),
	})

	if err != nil {
		// Handle duplicate key error (vehicle already parked)
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, &genericresponse.GenericResponse{
				StatusCode: http.StatusBadRequest,
				Message:    "Vehicle already in parking space",
			}
		}
		// Handle other internal errors
		return nil, &genericresponse.GenericResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	// Prepare the response with parking ticket information
	resp := &model.ParkVehicleResponse{
		ParkingTicket: model.ParkingTicket{
			VehicleNumber: req.VehicleNumber,
			VehicleID:     int(req.VehicleID),
			EntryTime:     time.Now(),
		},
	}

	// Set the parking lot name based on the ParkingLotID
	switch req.ParkingLotID {
	case 1:
		resp.ParkingTicket.ParkingLot = "Parking Lot A"
	case 2:
		resp.ParkingTicket.ParkingLot = "Parking Lot B"
	}

	return resp, nil
}
