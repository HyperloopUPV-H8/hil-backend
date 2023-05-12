package utilities

import "math/rand"

type VehicleState struct {
	XDistance   float64 `json:"xDistance"` //It goes between 22mm and 10 mm
	Current     float64 `json:"current"`
	Duty        byte    `json:"duty"`
	Temperature float64 `json:"temperature"`
}

func CreateVehicleState() VehicleState { //TODO
	VehicleState := &VehicleState{}
	VehicleState.XDistance = float64(rand.Intn(12) + 10)
	VehicleState.Current = float64(rand.Intn(20))
	VehicleState.Duty = byte(rand.Intn(100))
	VehicleState.Temperature = float64(rand.Intn(40) + 20)
	return *VehicleState
}
