package main

import (
	"encoding/binary"
	"errors"
	"fmt"
)

const VEHICLE_STATE_ID = 1
const FRONT_ORDER_ID = 2
const CONTROL_ORDER_ID = 3

//type HilParser map[int]any

// CreateHilParser() {
// 	return HilParser(
// 		map[int]Container
// 	)
// }

// type Container interface {
// 	Encode
// 	Decode
// }

// func (parser *HilParser) Encode() {}
// func (parser *HilParser) Decode() {}

// func CreateHilParser() *HilParser {
// 	return &HilParser{} //FIXME: Is it necessary?
// }

func Encode(data interface{}) []byte {
	switch dataType := data.(type) {
	case []VehicleState:
		head := make([]byte, 2)
		binary.LittleEndian.PutUint16(head, VEHICLE_STATE_ID)
		return Prepend(GetAllBytesFromVehiclesState(dataType), head...) //FIXME: Check prefix
	case []FrontOrder: //TODO: Is is necessary to diferenciate it? ControlOrder and FrontOrder? Don't think so
		head := make([]byte, 2)
		binary.LittleEndian.PutUint16(head, FRONT_ORDER_ID)
		return Prepend(GetAllBytesFromOrder(dataType), head...)
	case []ControlOrder:
		head := make([]byte, 2)
		binary.LittleEndian.PutUint16(head, CONTROL_ORDER_ID)
		return Prepend(GetAllBytesFromOrder(dataType), head...)
	default:
		fmt.Println("Does NOT match any type")
		return nil
	}
}

func Decode(data []byte) (any, error) { //FIXME: With a map choose the struct, define how to know it

	dataType := binary.LittleEndian.Uint64(data[0:2]) //FIXME: Little Endian?
	switch dataType {
	case VEHICLE_STATE_ID: //FIXME: Talk about types
		vehicleStates, err := GetAllVehicleStates(data[2:])
		return vehicleStates, err
	case FRONT_ORDER_ID: //TODO
		return nil, nil
	case CONTROL_ORDER_ID: //TODO
		return nil, nil
	default:
		fmt.Println("Does NOT match any type")
		return nil, errors.New("Does NOT match any type")
	}

}

func Decode1(data []byte) (interface{}, error) { //FIXME: With a map choose the struct, define how to know it
	vehicleStates, err := GetAllVehicleStates(data[:])
	return vehicleStates, err
}
