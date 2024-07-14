package service

import (
	"context"
	"parking_lot_service/internal/repo"
	"parking_lot_service/internal/service/model"
)

type ParkingLotService interface {
	GetFreeParkingSpaces(ctx context.Context) ([]*model.FreeSpotsResponse, error)
	GetFreeParkingSpaceById(ctx context.Context, parkingLotId int) (*model.FreeSpotsResponse, error)
	ParkVehicle(ctx context.Context, req *model.ParkVehicleRequest) (*model.ParkVehicleResponse, error)
	UnParkVehicle(ctx context.Context, req *model.UnParkVehicleRequest) (*model.UnParkVehicleResponse, error)
}

type impl struct {
	parkingLotRepo repo.ParkingLotRepo
}

func NewParkingLotService(parkingLotRepo repo.ParkingLotRepo) ParkingLotService {
	return &impl{parkingLotRepo: parkingLotRepo}
}
