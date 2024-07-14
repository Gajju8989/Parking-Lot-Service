package migration

import (
	"gorm.io/gorm"
	"parking_lot_service/internal/repo/models"
)

func MigrateAll(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.ParkingSpace{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&models.ParkedVehicle{}); err != nil {
		return err
	}
	return nil
}
