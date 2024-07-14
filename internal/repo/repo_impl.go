package repo

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"parking_lot_service/internal/repo/models"
)

func (s *impl) SeedParkingSpace(ctx context.Context) error {
	var count int64
	err := s.db.
		WithContext(ctx).
		Model(&models.ParkingSpace{}).
		Count(&count).
		Error

	if err != nil {
		return fmt.Errorf("error counting existing parking spaces: %w", err)
	}

	if count > 0 {
		return nil // Skip seeding if records already exist
	}

	parkingSpaces := []models.ParkingSpace{
		// Parking Lot A
		{ParkingLotId: models.ParkingLotA, VehicleTypeId: models.MotorcyclesAndScooters, AvailableSpots: 50},
		{ParkingLotId: models.ParkingLotA, VehicleTypeId: models.CarsAndSUVs, AvailableSpots: 30},
		{ParkingLotId: models.ParkingLotA, VehicleTypeId: models.BusesAndTrucks, AvailableSpots: 20},
		// Parking Lot B
		{ParkingLotId: models.ParkingLotB, VehicleTypeId: models.MotorcyclesAndScooters, AvailableSpots: 100},
		{ParkingLotId: models.ParkingLotB, VehicleTypeId: models.CarsAndSUVs, AvailableSpots: 80},
		{ParkingLotId: models.ParkingLotB, VehicleTypeId: models.BusesAndTrucks, AvailableSpots: 40},
	}

	// Insert all parking spaces in a single transaction
	tx := s.db.WithContext(ctx).Begin()
	for _, space := range parkingSpaces {
		err = tx.
			Create(&space).
			Error
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error creating parking space: %w", err)
		}
	}

	err = tx.
		Commit().
		Error

	if err != nil {
		return fmt.Errorf("error committing parking spaces: %w", err)
	}

	return nil
}

func (s *impl) GetParkingSpaces(ctx context.Context) ([]*models.ParkingSpace, error) {
	var parkingSpaces []*models.ParkingSpace

	// Query the database to fetch all parking spaces
	res := s.db.
		WithContext(ctx).
		Find(&parkingSpaces)

	if res.Error != nil {
		return nil, res.Error
	}

	return parkingSpaces, nil
}

func (s *impl) GetFreeParkingSpaceById(ctx context.Context, parkingLotId int) ([]*models.ParkingSpace, error) {
	var parkingSpaces []*models.ParkingSpace

	// Query parking spaces where ParkingLotId matches
	err := s.db.WithContext(ctx).
		Where("parking_lot_id = ?", parkingLotId).
		Find(&parkingSpaces).
		Error

	if err != nil {
		return nil, err
	}

	return parkingSpaces, nil
}

func (s *impl) DecreaseAvailableSpot(ctx context.Context, parkingSpace *models.ParkingSpace) error {
	result := s.
		db.
		WithContext(ctx).
		Model(&models.ParkingSpace{}).
		Where("parking_lot_id = ?", parkingSpace.ParkingLotId).
		Where("vehicle_type_id = ?", parkingSpace.VehicleTypeId).
		UpdateColumn("available_spots", gorm.Expr("available_spots - ?", 1))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("no parking space updated")
	}

	return nil
}

func (s *impl) SaveParkedVehicle(ctx context.Context, parkingSpace *models.ParkedVehicle) error {
	err := s.
		db.
		WithContext(ctx).
		Create(&parkingSpace).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (s *impl) IncreaseAvailableSpot(ctx context.Context, parkingSpace *models.ParkedVehicle) error {
	return nil
}