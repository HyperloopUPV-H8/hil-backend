package utilities

type VehicleState struct {
	XDistance   float64 `json:"xDistance"` //It goes between 22mm and 10 mm
	Current     float64 `json:"current"`
	Duty        byte    `json:"duty"`
	Temperature float64 `json:"temperature"`
}

type Perturbation struct {
	Id    string `json:"id"`
	TypeP string `json:"type"`
	Value int    `json:"value"`
}

type PerturbationArray []Perturbation
