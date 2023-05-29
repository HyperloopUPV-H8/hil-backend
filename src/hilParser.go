package main

import (
	"encoding/binary"
	"errors"
	"fmt"
)

const VEHICLE_STATE_ID = 1
const FRONT_ORDER_ID = 2
const CONTROL_ORDER_ID = 3

func Encode(data interface{}) []byte {
	switch dataType := data.(type) {
	case []VehicleState:
		head := make([]byte, 2)
		binary.LittleEndian.PutUint16(head, VEHICLE_STATE_ID)
		return Prepend(GetAllBytesFromVehiclesState(dataType), head...)
	case []Order:
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

	dataType := binary.LittleEndian.Uint16(data[0:2])
	switch dataType {
	case VEHICLE_STATE_ID:
		return GetAllVehicleStates(data[2:])
	case FRONT_ORDER_ID: //TODO
		return nil, nil
	case CONTROL_ORDER_ID: //TODO
		return GetAllControlOrders(data[2:])
	default:
		fmt.Println("Does NOT match any type (decode)")
		return nil, errors.New("Does NOT match any type")
	}

}
