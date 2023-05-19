package main

import (
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
		return CreateVehicleStateFromBytes(dataType)
	default:
		fmt.Println("Does NOT match any type")
		return nil
	}
}

func Decode(dataType string, data []byte) any { //FIXME: With a map choose the struct, define how to know it
	switch dataType {
	case "VehicleState":
		return GetVehicleState(data)
	default:
		fmt.Println("Does NOT match any type")
		return nil
	}

}
