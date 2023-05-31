package main

import (
	"encoding/binary"
	"errors"
	"main/conversions"
	"main/models"

	trace "github.com/rs/zerolog/log"
)

const VEHICLE_STATE_ID = 1
const FRONT_ORDER_ID = 2
const CONTROL_ORDER_ID = 3

func Encode(data interface{}) []byte { //FIXME: For encode array of structs, order.Bytes implemented for adding always the prefix
	switch dataType := data.(type) {
	case []models.VehicleState:
		head := make([]byte, 2)
		binary.LittleEndian.PutUint16(head, VEHICLE_STATE_ID)
		return Prepend(conversions.GetAllBytesFromVehiclesState(dataType), head...)
	case []models.Order:
		head := make([]byte, 2)
		switch dataType[0].(type) {
		case models.FrontOrder:
			binary.LittleEndian.PutUint16(head, FRONT_ORDER_ID)
		case models.ControlOrder:
			binary.LittleEndian.PutUint16(head, CONTROL_ORDER_ID)
		default:
			trace.Warn().Msg("Does NOT match any ORDER type")
			return nil
		}

		return Prepend(conversions.GetAllBytesFromOrder(dataType), head...)
	default:
		trace.Warn().Msg("Does NOT match any type")
		return nil
	}
}

func Decode(data []byte) (any, error) {

	dataType := binary.LittleEndian.Uint16(data[0:2])
	switch dataType {
	case VEHICLE_STATE_ID:
		return conversions.GetAllVehicleStates(data[2:])
	case FRONT_ORDER_ID: //TODO: Not needed now
		return nil, errors.New("Decode for this type of order not implemented")
	case CONTROL_ORDER_ID:
		return conversions.GetAllControlOrders(data[2:])
	default:
		trace.Warn().Msg("Does NOT match any type")
		return nil, errors.New("does not match any type in decode function")
	}

}
