package models

type VehicleState struct {
	YDistance   float64 `json:"yDistance"` //Value between 22mm and 10 mm
	Current     float64 `json:"current"`
	Duty        byte    `json:"duty"`
	Temperature float64 `json:"temperature"`
}

type Perturbation struct {
	Id    string `json:"id"`
	TypeP string `json:"type"`
	Value int    `json:"value"`
}

type PerturbationOrder []Perturbation

type InputData struct {
	Id       string   `json:"id"`
	Type     string   `json:"type"`
	Value    float64  `json:"value"`
	Enabled  bool     `json:"enabled"`
	Validity Validity `json:"validity"`
}

type Validity struct {
	IsValid bool   `json:"isValid"`
	Msg     string `json:"msg"`
}

type FormData []InputData
