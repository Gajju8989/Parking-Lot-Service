package repo

import (
	"context"
	"fmt"
	"parking_lot_service/internal/repo/models"
)

// SeedParkingSpace seeds the database with initial parking space records if none exist.
// It checks for existing records, and if none are found, it inserts predefined parking spaces.
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

// GetParkingSpaces retrieves all parking spaces from the database.
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

// GetFreeParkingSpaceById retrieves free parking spaces for a given parking lot ID.
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

// SaveParkedVehicle saves a parked vehicle record to the database.
func (s *impl) SaveParkedVehicle(ctx context.Context, vehicleDetail *models.ParkedVehicle) error {
	err := s.
		db.
		WithContext(ctx).
		Create(&vehicleDetail).
		Error
	if err != nil {
		return err
	}
	return nil
}

// GetAvailableParkingSpotsByParkingLotIdAndVehicleId counts available parking spots by parking lot ID and vehicle type ID.
func (s *impl) GetAvailableParkingSpotsByParkingLotIdAndVehicleId(ctx context.Context,
	parkingLotId, vehicleId int) (int, error) {
	resp := models.ParkingSpace{}
	err := s.db.WithContext(ctx).
		Model(&models.ParkingSpace{}).
		Where("parking_lot_id = ? AND vehicle_type_id = ?", parkingLotId, vehicleId).
		First(&resp).
		Error

	if err != nil {
		return 0, err
	}
	return resp.AvailableSpots, nil
}

// UpdateParkingSpace updates the available spots of a parking space in the database.
func (s *impl) UpdateParkingSpace(ctx context.Context, parkingSpace *models.ParkingSpace) error {
	err := s.
		db.
		WithContext(ctx).
		Model(&models.ParkingSpace{}).
		Where("parking_lot_id = ? AND vehicle_type_id = ?", parkingSpace.ParkingLotId, parkingSpace.VehicleTypeId).
		Updates(map[string]interface{}{
			"available_spots": parkingSpace.AvailableSpots,
		}).
		Error

	if err != nil {
		return err
	}
	return nil
}

// GetParkedVehicle checks if a parked vehicle exists in the database and returns it.
func (s *impl) GetParkedVehicle(ctx context.Context, vehicleNumber string) (*models.ParkedVehicle, error) {
	var existingVehicle models.ParkedVehicle

	// Check if the vehicle exists
	err := s.db.WithContext(ctx).
		Where("vehicle_number = ?", vehicleNumber).
		First(&existingVehicle).
		Error

	if err != nil {
		return nil, err
	}

	return &existingVehicle, nil
}

// DeleteParkedVehicle deletes the specified parked vehicle from the database.
func (s *impl) DeleteParkedVehicle(ctx context.Context, parkedVehicle *models.ParkedVehicle) error {
	// Delete the vehicle
	err := s.db.WithContext(ctx).
		Delete(parkedVehicle).
		Error

	if err != nil {
		return err
	}

	return nil
}
