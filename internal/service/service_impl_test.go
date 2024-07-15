package service

import (
	"testing"
	"time"
)

func Test_calculateFare(t *testing.T) {
	type args struct {
		parkingLotID  int
		vehicleTypeId int
		duration      time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "Valid case for Motorcycles/scooters at Parking Lot 1",
			args: args{
				parkingLotID:  1,
				vehicleTypeId: 1,
				duration:      2 * time.Hour,
			},
			want:    10.0, // 2 hours * HourlyRate: 5
			wantErr: false,
		},
		{
			name: "Case Where MotorCycles/Scooters at Parking Lot 1 with some extra minutes",
			args: args{
				parkingLotID:  1,
				vehicleTypeId: 1,
				duration:      2*time.Hour + 10*time.Minute,
			},
			want:    15.0, // 3 hours * HourlyRate: 5
			wantErr: false,
		},
		{
			name: "Valid case for Cars/SUVs at Parking Lot 1",
			args: args{
				parkingLotID:  1,
				vehicleTypeId: 2,
				duration:      28 * time.Hour,
			},
			want:    574.0, // 28 hours * HourlyRate: 20.50/hr
			wantErr: false,
		},
		{
			name: "Valid case for Cars/SUVs at Parking Lot 1 with Extra Minutes",
			args: args{
				parkingLotID:  1,
				vehicleTypeId: 2,
				duration:      28*time.Hour + 1*time.Minute,
			},
			want:    594.5, // 29 hours * HourlyRate: 20.50/hr
			wantErr: false,
		},
		{
			name: "Valid case for Bus/Trucks at Parking Lot 1",
			args: args{
				parkingLotID:  1,
				vehicleTypeId: 3,
				duration:      23 * time.Hour,
			},
			want:    1150.0, // 23 hours * HourlyRate: 50
			wantErr: false,
		},
		{
			name: "Valid case for Bus/Trucks at Parking Lot 1 with 2 days 2 hour and 1 minute",
			args: args{
				parkingLotID:  1,
				vehicleTypeId: 3,
				duration:      50*time.Hour + 1*time.Minute,
			},
			want:    1850.0, // for first day=1200 based on 24*50 ,then second day 500, and for 2 hour 1 minutes 150.
			wantErr: false,
		},
		{
			name: "Invalid case: No tariff found",
			args: args{
				parkingLotID:  3, // Non-existing parking lot ID
				vehicleTypeId: 1,
				duration:      2 * time.Hour,
			},
			want:    0.0,
			wantErr: true,
		},
		{
			name: "Valid case for Motorcycles/scooters at Parking Lot 2",
			args: args{
				parkingLotID:  2,
				vehicleTypeId: 1,
				duration:      2 * time.Hour,
			},
			want:    21.0, // 2 hours * HourlyRate: 10.5
			wantErr: false,
		},
		{
			name: "Valid case for Motorcycles/scooters at Parking Lot 2 with extra minutes",
			args: args{
				parkingLotID:  2,
				vehicleTypeId: 1,
				duration:      2*time.Hour + 10*time.Minute,
			},
			want:    31.5, // 3 hours * HourlyRate: 10.5
			wantErr: false,
		},
		{
			name: "Valid case for Cars/SUVs at Parking Lot 2",
			args: args{
				parkingLotID:  2,
				vehicleTypeId: 2,
				duration:      2 * time.Hour,
			},
			want:    75.0, // FirstHourRate: 50 + AdditionalHourRate: 25
			wantErr: false,
		},
		{
			name: "Valid case for Cars/SUVs at Parking Lot 2 with extra minutes",
			args: args{
				parkingLotID:  2,
				vehicleTypeId: 2,
				duration:      2*time.Hour + 10*time.Minute,
			},
			want:    100.0, //FirstHourRate: 50 + 1 AdditionalHourRate: 25+10 Extra Minutes mean 1 additional Hour (25)
			wantErr: false,
		},
		{
			name: "Valid case for Cars/SUVs at Parking Lot 2 with extra hours without extra minutes",
			args: args{
				parkingLotID:  2,
				vehicleTypeId: 2,
				duration:      2 * time.Hour,
			},
			want:    75.0, //FirstHourRate: 50 + 1 AdditionalHourRate: 25+10 Extra Minutes mean 1 additional Hour (25)
			wantErr: false,
		},
		{
			name: "Valid case for Cars/SUVs at Parking Lot 2 with 5 hours",
			args: args{
				parkingLotID:  2,
				vehicleTypeId: 2,
				duration:      5 * time.Hour,
			},
			want:    150.0, // FirstHourRate: 50 + 4 * AdditionalHourRate: 25
			wantErr: false,
		},
		{
			name: "Valid case for Bus/Trucks at Parking Lot 2",
			args: args{
				parkingLotID:  2,
				vehicleTypeId: 3,
				duration:      2 * time.Hour,
			},
			want:    200.0, // 2 hours * HourlyRate: 100
			wantErr: false,
		},
		{
			name: "Valid case for Bus/Trucks at Parking Lot 2 with extra minutes",
			args: args{
				parkingLotID:  2,
				vehicleTypeId: 3,
				duration:      2*time.Hour + 10*time.Minute,
			},
			want:    300.0, // 3 hours * HourlyRate: 100
			wantErr: false,
		},
		{
			name: "Edge case: Duration exactly a day for Parking Lot 1 Buses/Trucks",
			args: args{
				parkingLotID:  1,
				vehicleTypeId: 3,
				duration:      24 * time.Hour,
			},
			want:    1200.0, // DayRate: 1000
			wantErr: false,
		},
		{
			name: "Edge case: Duration just above a day for Parking Lot 1 Buses/Trucks",
			args: args{
				parkingLotID:  1,
				vehicleTypeId: 3,
				duration:      24*time.Hour + 1*time.Minute,
			},
			want:    1250.0, // DayRate: 24*50  + 50 1 Minute
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculateFare(tt.args.parkingLotID, tt.args.vehicleTypeId, tt.args.duration)
			if (err != nil) != tt.wantErr {
				t.Errorf("calculateFare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got != tt.want {
				t.Errorf("calculateFare() got = %v, want %v", got, tt.want)
			}
		})
	}
}
