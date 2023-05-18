package utilities

import (
	"math"
	"math/rand"
)

func RandomVehicleState() VehicleState {
	VehicleState := &VehicleState{}
	VehicleState.YDistance = float64(rand.Intn(13)+10) + (math.Round(rand.Float64()*100) / 100)
	VehicleState.Current = float64(rand.Intn(20)) + (math.Round(rand.Float64()*100) / 100)
	VehicleState.Duty = byte(rand.Intn(100))
	VehicleState.Temperature = float64(rand.Intn(40)+20) + (math.Round(rand.Float64()*100) / 100)
	return *VehicleState
}
