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
	ParkingLotID  models.ParkingLot  `json:"parking_lot_id" binding:"required"`
	VehicleNumber string             `json:"vehicle_number" binding:"required"`
	VehicleID     models.VehicleType `json:"vehicle_id" binding:"required"`
}

type UnParkVehicleResponse struct {
	Parking ParkingReceipt `json:"parking_receipt"`
}

type ParkingReceipt struct {
	VehicleNumber string  `json:"vehicle_number"`
	TotalFare     float64 `json:"total_fare"`
	From          string  `json:"from"`
	To            string  `json:"to"`
	VehicleID     int     `json:"vehicle_id"`
	ParkingLotID  int     `json:"parking_lot_id"`
}

// Tariff represents the tariff details for different vehicle types in a parking lot.
type Tariff struct {
	HourlyRate            float64
	DayRate               float64
	FirstHourRate         float64
	AdditionalHourRate    float64
	MaxDurationForDayRate time.Duration
}
