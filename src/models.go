package main

type VehicleState struct {
	YDistance   float64 `json:"yDistance"` //It goes between 22mm and 10 mm
	Current     float64 `json:"current"`
	Duty        byte    `json:"duty"`
	Temperature float64 `json:"temperature"`
}

type Perturbation struct {
	Id    string `json:"id"`
	TypeP string `json:"type"`
	Value int    `json:"value"`
}

// type SimulationData struct {
// 	Current  float64
// 	Distance float64
// }

type FrontOrder struct {
	kind    string
	payload []byte
}

type ControlOrder struct {
	Variable string
	State    bool
}

type PerturbationOrder []Perturbation

type Order interface {
	Bytes() []byte
	Read([]byte)
	//GetAllBytesFromOrder() []byte
}
