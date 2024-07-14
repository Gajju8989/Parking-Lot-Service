package model

import (
	"parking_lot_service/internal/repo/models"
	"time"
)

type FreeSpotsResponse struct {
	ParkingLotID                    int `json:"parkingLotId"`
	FreeSpotsForMotorcyclesScooters int `json:"freeSpotsForMotorcyclesScooters"`
	FreeSpotsForCarsSUVs            int `json:"freeSpotsForCarsSUVs"`
	FreeSpotsForBusesTrucks         int `json:"freeSpotsForBusesTrucks"`
}

type ParkVehicleRequest struct {
	ParkingLotID  models.ParkingLot  `json:"parking_lot_id" binding:"required"`
	VehicleID     models.VehicleType `json:"vehicle_id" binding:"required"`
	VehicleNumber string             `json:"vehicle_number" binding:"required"`
	VehicleName   string             `json:"vehicle_name"`
}

type ParkVehicleResponse struct {
	ParkingTicket ParkingTicket `json:"parking_ticket"`
}
type ParkingTicket struct {
	VehicleNumber string    `json:"vehicle_number"`
	ParkingLot    string    `json:"parking_lot"`
	VehicleID     int       `json:"vehicle_id"`
	EntryTime     time.Time `json:"entry_time"`
}

type UnParkVehicleRequest struct {
	ParkingLotID  models.ParkingLot `json:"parking_lot_id"`
	VehicleNumber string            `json:"vehicle_number"`
	VehicleID     int               `json:"vehicle_id"`
}

type UnParkVehicleResponse struct {
	Parking ParkingReceipt `json:"parking_receipt"`
}
type ParkingReceipt struct {
}
