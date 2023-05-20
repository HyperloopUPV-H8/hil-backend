package main

import (
	"errors"
	"fmt"
)

type HilParser map[int]any

// CreateHilParser() {
// 	return HilParser(
// 		map[int]Container
// 	)
// }

// type Container interface {
// 	Encode
// 	Decode
// }

//func (parser *HilParser) Encode() {}
//func (parser *HilParser) Decode() {}

func CreateHilParser() *HilParser {
	return &HilParser{} //FIXME
}

func Encode(data interface{}) []byte {
	switch dataType := data.(type) {
	case VehicleState:
		return CreateBytesFromVehicleState(dataType) //TODO: add prefix of type of message
	default:
		fmt.Println("Does NOT match any type")
		return nil
	}
}

func Decode(data []byte) (any, error) { //FIXME: With a map choose the struct, define how to know it
	dataType := string(data[0:2])
	switch dataType {
	case "VehicleState": //FIXME: Talk about types
		vehicleStates, err := GetAllVehicleStates(data[2:])
		return vehicleStates, err
	default:
		fmt.Println("Does NOT match any type")
		return nil, errors.New("Does NOT match any type")
	}

}
