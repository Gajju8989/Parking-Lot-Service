package repo

import (
	"context"
	"gorm.io/gorm"
	"parking_lot_service/internal/repo/models"
)

type ParkingLotRepo interface {
	SeedParkingSpace(ctx context.Context) error
	GetParkingSpaces(ctx context.Context) ([]*models.ParkingSpace, error)
	GetFreeParkingSpaceById(ctx context.Context, parkingLotId int) ([]*models.ParkingSpace, error)
	SaveParkedVehicle(ctx context.Context, vehicleDetail *models.ParkedVehicle) error
	GetAvailableParkingSpotsByParkingLotIdAndVehicleId(ctx context.Context, parkingLotId, VehicleId int) (int, error)
	UpdateParkingSpace(ctx context.Context, parkingSpace *models.ParkingSpace) error
	GetParkedVehicle(ctx context.Context, vehicleNumber string) (*models.ParkedVehicle, error)
	DeleteParkedVehicle(ctx context.Context, parkedVehicle *models.ParkedVehicle) error
}

type impl struct {
	db *gorm.DB
}

func NewParkingLotRepo(db *gorm.DB) ParkingLotRepo {
	return &impl{db: db}
}
