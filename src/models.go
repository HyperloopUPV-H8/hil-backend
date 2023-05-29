package main

import (
	"bytes"
	"encoding/binary"
)

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

type PerturbationOrder []Perturbation

// type SimulationData struct {
// 	Current  float64
// 	Distance float64
// }

type Order interface {
	Bytes() []byte
	//Read([]byte) TODO
}

type FrontOrder struct {
	Kind    string  `json:"kind"`
	Payload float64 `json:"payload"`
}

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

func (order FrontOrder) Bytes() []byte {
	buf1 := []byte(order.Kind)
	var buf2 [8]byte
	binary.LittleEndian.PutUint64(buf2[:], uint64(order.Payload))
	return append(buf1, buf2[:]...)
}

func (order *FrontOrder) Read(data []byte) {
	reader := bytes.NewReader(data)
	binary.Read(reader, binary.LittleEndian, order)
}

type ControlOrder struct {
	Id    uint8 `json:"id"`
	State bool  `json:"state"`
}

func (order ControlOrder) Bytes() []byte {
	//var buf1 [8]byte
	//binary.LittleEndian.PutUint64(buf1[:], order.Id)
	buf1 := []byte{order.Id}
	var booleanValue uint8
	if order.State {
		booleanValue = 1
	} else {
		booleanValue = 0
	}
	return append(buf1, booleanValue)
}

func (order *ControlOrder) Read(data []byte) {
	reader := bytes.NewReader(data)
	binary.Read(reader, binary.LittleEndian, order)
}
