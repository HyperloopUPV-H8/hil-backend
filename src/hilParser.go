package main

import (
	"encoding/binary"
	"errors"
	"fmt"
)

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
	case VehicleState:
		return CreateBytesFromVehicleState(dataType) //TODO: add prefix of type of message
	case Order: //FIXME: Is is necessary to diferenciate it? COntrolOrder and FrontOrder?
		return dataType.Bytes() //TODO: add prefix of type of message
	default:
		fmt.Println("Does NOT match any type")
		return nil
	}
}

func Decode(data []byte) (any, error) { //FIXME: With a map choose the struct, define how to know it
	dataType := binary.LittleEndian.Uint64(data[0:2]) //FIXME: Little Endian?
	switch dataType {
	case 1: //FIXME: Talk about types
		vehicleStates, err := GetAllVehicleStates(data[2:])
		return vehicleStates, err
	default:
		fmt.Println("Does NOT match any type")
		return nil, errors.New("Does NOT match any type")
	}

}

func Decode1(data []byte) (interface{}, error) { //FIXME: With a map choose the struct, define how to know it
	vehicleStates, err := GetAllVehicleStates(data[:])
	return vehicleStates, err
}
