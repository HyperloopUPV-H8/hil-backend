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
	VehicleState.XDistance = float64(rand.Intn(12) + 10.0)
	//VehicleState.XRotation = 0
	return *VehicleState
}
