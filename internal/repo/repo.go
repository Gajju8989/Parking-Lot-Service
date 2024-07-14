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
	DecreaseAvailableSpot(ctx context.Context, parkingSpace *models.ParkingSpace) error
	SaveParkedVehicle(ctx context.Context, parkingSpace *models.ParkedVehicle) error
	IncreaseAvailableSpot(ctx context.Context, parkingSpace *models.ParkedVehicle) error
}

type impl struct {
	db *gorm.DB
}

func NewParkingLotRepo(db *gorm.DB) ParkingLotRepo {
	return &impl{db: db}
}
