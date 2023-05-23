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
	case []Order: //TODO: Is is necessary to diferenciate it? ControlOrder and FrontOrder? Don't think so
		head := make([]byte, 2)
		switch dataType[0].(type) { //FIXME, is it correct?
		case FrontOrder:
			binary.LittleEndian.PutUint16(head, FRONT_ORDER_ID)
		case ControlOrder:
			binary.LittleEndian.PutUint16(head, CONTROL_ORDER_ID)
		default:
			fmt.Println("Does NOT match any ORDER type (Encode)")
			return nil
		}

		return Prepend(GetAllBytesFromOrder(dataType), head...)
	default:
		fmt.Println("Does NOT match any type (Encode2)")
		return nil
	}
}

func Decode(data []byte) (any, error) { //FIXME: With a map choose the struct, define how to know it

	dataType := binary.LittleEndian.Uint16(data[0:2]) //FIXME: Little Endian?
	fmt.Println("Datatype Decode: ", dataType)
	switch dataType {
	case VEHICLE_STATE_ID: //FIXME: Talk about types
		vehicleStates, err := GetAllVehicleStates(data[2:])
		return vehicleStates, err
	case FRONT_ORDER_ID: //TODO
		return nil, nil
	case CONTROL_ORDER_ID: //TODO
		return nil, nil
	default:
		fmt.Println("Does NOT match any type (decode)")
		return nil, errors.New("Does NOT match any type")
	}

}

func Decode1(data []byte) (interface{}, error) { //FIXME: With a map choose the struct, define how to know it
	vehicleStates, err := GetAllVehicleStates(data[:])
	return vehicleStates, err
}
