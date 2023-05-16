package utilities

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

func ConvertBytesToFloat64(bytes [8]byte) float64 {
	num := math.Float64frombits(binary.LittleEndian.Uint64(bytes[:]))
	return num
}

func CreateMockBytes(vehicleState VehicleState) [25]byte {

	buf1 := ConvertFloat64ToBytes(vehicleState.XDistance)
	buf2 := ConvertFloat64ToBytes(vehicleState.Current)
	var buf3 [1]byte = [1]byte{vehicleState.Duty}
	buf4 := ConvertFloat64ToBytes(vehicleState.Temperature)

	result := append(append(append(buf1[:], buf2[:]...), buf3[:]...), buf4[:]...)

	var mockBytes [25]byte
	copy(mockBytes[:], result)
	return mockBytes
}

func TestGetVehicleState(vehicleState VehicleState) VehicleState {
	bytes := CreateMockBytes(vehicleState)
	return GetVehicleState(bytes[:])

}
