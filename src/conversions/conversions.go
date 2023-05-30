package conversions

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"main/models"
	"math"
)

const VEHICLE_STATE_LENGTH = 25

func GetVehicleState(data []byte) models.VehicleState {
	reader := bytes.NewReader(data)
	vehicleState := &models.VehicleState{}
	binary.Read(reader, binary.LittleEndian, vehicleState)
	return *vehicleState
}

func GetAllVehicleStates(data []byte) ([]models.VehicleState, error) {
	vehicleStateArray := []models.VehicleState{}
	reader := bytes.NewReader(data)
	var err error
	for i := 0; i <= len(data)-VEHICLE_STATE_LENGTH; i += VEHICLE_STATE_LENGTH {
		vehicleState := &models.VehicleState{}
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

func GetBytesFromVehicleState(vehicleState models.VehicleState) []byte {

	buf1 := ConvertFloat64ToBytes(vehicleState.YDistance)
	buf2 := ConvertFloat64ToBytes(vehicleState.Current)
	var buf3 [1]byte = [1]byte{vehicleState.Duty}
	buf4 := ConvertFloat64ToBytes(vehicleState.Temperature)

	return append(append(append(buf1[:], buf2[:]...), buf3[:]...), buf4[:]...)
}

func GetAllBytesFromVehiclesState(vehiclesState []models.VehicleState) []byte {
	var result []byte
	for _, vehicle := range vehiclesState {
		result = append(result, GetBytesFromVehicleState(vehicle)...)
	}
	return result
}

func GetAllControlOrders(data []byte) ([]models.ControlOrder, error) {
	ordersArray := []models.ControlOrder{}
	reader := bytes.NewReader(data)
	var err error
	for reader.Len() > 0 {
		order := &models.ControlOrder{}
		err = binary.Read(reader, binary.LittleEndian, order)
		if err != nil {
			fmt.Println("error decoding control orders: ", err)
			break
		}
		ordersArray = append(ordersArray, *order)
	}
	return ordersArray, err
}

func GetAllBytesFromOrder(data []models.Order) []byte {
	var result []byte
	for _, order := range data {
		result = append(result, order.Bytes()...)
	}
	return result
}

func GetAllBytesFromControlOrder(data []models.ControlOrder) []byte {
	var result []byte
	for _, order := range data {
		result = append(result, order.Bytes()...)
	}
	return result
}

func ConvertFormDataToOrders(form models.FormData) []models.FrontOrder {
	var frontOrders []models.FrontOrder
	for _, order := range form {
		if order.Enabled && order.Validity.IsValid { //FIXME: Check type, not necessary now
			frontOrder := models.FrontOrder{Kind: order.Id, Payload: order.Value}
			frontOrders = append(frontOrders, frontOrder)
		}
	}
	return frontOrders
}
