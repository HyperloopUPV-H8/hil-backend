package main

import (
	"bytes"
	"encoding/binary"
	"math"
)

func GetVehicleState(data []byte) VehicleState {
	reader := bytes.NewReader(data)
	vehicleState := &VehicleState{}
	binary.Read(reader, binary.LittleEndian, vehicleState)
	return *vehicleState
}

func ConvertFloat64ToBytes(num float64) [8]byte {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], math.Float64bits(num))
	return buf
}

func CreateVehicleStateFromBytes(vehicleState VehicleState) []byte {

	buf1 := ConvertFloat64ToBytes(vehicleState.YDistance)
	buf2 := ConvertFloat64ToBytes(vehicleState.Current)
	var buf3 [1]byte = [1]byte{vehicleState.Duty}
	buf4 := ConvertFloat64ToBytes(vehicleState.Temperature)

	result := append(append(append(buf1[:], buf2[:]...), buf3[:]...), buf4[:]...)

	return result
}
