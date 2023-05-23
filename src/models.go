package main

import (
	"bytes"
	"encoding/binary"
	"strconv"
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
	Read([]byte)
	//GetAllBytesFromOrder() []byte
}

type FrontOrder struct {
	kind    string
	payload []byte
}

func (order FrontOrder) Bytes() []byte {
	buf1 := []byte(order.kind)
	return append(buf1, order.payload...)
}

func (order FrontOrder) Read(data []byte) { //TODO: Add prefix
	reader := bytes.NewReader(data)
	binary.Read(reader, binary.LittleEndian, order)
}

type ControlOrder struct {
	Variable string
	State    bool
}

func (order ControlOrder) Bytes() []byte {
	buf1 := []byte(order.Variable)
	return strconv.AppendBool(buf1, order.State)
}

func (order ControlOrder) Read(data []byte) { //TODO: Add prefix
	reader := bytes.NewReader(data)
	binary.Read(reader, binary.LittleEndian, order)
}
