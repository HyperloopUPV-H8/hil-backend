package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

const VEHICLE_STATE_LENGTH = 25

func GetVehicleState(data []byte) VehicleState {
	reader := bytes.NewReader(data)
	vehicleState := &VehicleState{}
	binary.Read(reader, binary.LittleEndian, vehicleState)
	return *vehicleState
}

func GetAllVehicleStates(data []byte) ([]VehicleState, error) {
	vehicleStateArray := []VehicleState{}
	reader := bytes.NewReader(data)
	var err error
	for i := 0; i <= len(data)-VEHICLE_STATE_LENGTH; i += VEHICLE_STATE_LENGTH {
		vehicleState := &VehicleState{}
		err = binary.Read(reader, binary.LittleEndian, vehicleState)
		if err != nil {
			break
		}
		vehicleStateArray = append(vehicleStateArray, *vehicleState)
	}
	return vehicleStateArray, err
}

func ConvertFloat64ToBytes(num float64) [8]byte {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], math.Float64bits(num))
	return buf
}

func GetBytesFromVehicleState(vehicleState VehicleState) []byte {

	buf1 := ConvertFloat64ToBytes(vehicleState.YDistance)
	buf2 := ConvertFloat64ToBytes(vehicleState.Current)
	var buf3 [1]byte = [1]byte{vehicleState.Duty}
	buf4 := ConvertFloat64ToBytes(vehicleState.Temperature)

	return append(append(append(buf1[:], buf2[:]...), buf3[:]...), buf4[:]...)
}

func GetAllBytesFromVehiclesState(vehiclesState []VehicleState) []byte {
	var result []byte
	for _, vehicle := range vehiclesState {
		result = append(result, GetBytesFromVehicleState(vehicle)...)
	}
	return result
}

func GetAllControlOrders(data []byte) ([]ControlOrder, error) {
	ordersArray := []ControlOrder{}
	reader := bytes.NewReader(data)
	var err error
	for reader.Len() > 0 { //FIXME?
		order := &ControlOrder{}
		err = binary.Read(reader, binary.LittleEndian, order) // TODO: There is an Error here
		if err != nil {
			fmt.Println("error decoding control orders: ", err)
			break
		}
		ordersArray = append(ordersArray, *order)
	}
	return ordersArray, err
}

func GetAllBytesFromOrder(data []Order) []byte {
	var result []byte
	for _, order := range data {
		result = append(result, order.Bytes()...)
	}
	return result
}

func GetAllBytesFromControlOrder(data []ControlOrder) []byte {
	var result []byte
	for _, order := range data {
		result = append(result, order.Bytes()...)
	}
	return result
}

func ConvertFormDataToOrders(form FormData) []FrontOrder {
	var frontOrders []FrontOrder
	for _, order := range form {
		if order.Enabled && order.Validity.IsValid { //TODO: Check type
			frontOrder := FrontOrder{Kind: order.Id, Payload: order.Value}
			fmt.Println(frontOrder)
			frontOrders = append(frontOrders, frontOrder)
		}
	}
	return frontOrders
}
