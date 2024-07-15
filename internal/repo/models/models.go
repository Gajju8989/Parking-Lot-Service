package models

import "time"

// ParkingLot represents different types of parking lots.
type ParkingLot int

const (
	ParkingLotA ParkingLot = 1
	ParkingLotB ParkingLot = 2
)

// VehicleType represents different types of vehicles.
type VehicleType int

const (
	MotorcyclesAndScooters VehicleType = 1
	CarsAndSUVs            VehicleType = 2
	BusesAndTrucks         VehicleType = 3
)

type ParkingSpace struct {
	ID             uint        `gorm:"primaryKey"` // Unique identifier for each parking space
	ParkingLotId   ParkingLot  `gorm:"not null;index:idx_parking_lot_vehicle_type"`
	VehicleTypeId  VehicleType `gorm:"not null;index:idx_parking_lot_vehicle_type"`
	AvailableSpots int         `gorm:"not null"` // Number of free spots left for the specified vehicle type
}

type ParkedVehicle struct {
	VehicleNumber string      `gorm:"primaryKey"`
	ParkingLotID  ParkingLot  `gorm:"not null"`
	VehicleTypeId VehicleType `gorm:"not null"`
	VehicleName   string      `gorm:"type:varchar(150)"`
	EntryTime     time.Time   `gorm:"not null"`
}
