package utilities

import (
	"testing"
)

type ConversionCase struct {
	buf    VehicleState
	result VehicleState
}

type ConversionCaseBytes struct {
	buf    []byte
	result VehicleState
}

func TestConversion(t *testing.T) {
	t.Run("struct conversion", func(t *testing.T) {
		cases := []ConversionCase{
			{buf: VehicleState{XDistance: 2.45, Current: 4.3, Duty: 1, Temperature: 10.2}, result: VehicleState{XDistance: 2.45, Current: 4.3, Duty: 1, Temperature: 10.2}},
			{buf: VehicleState{XDistance: 22.5, Current: 42.3, Duty: 91, Temperature: 0.2}, result: VehicleState{XDistance: 22.5, Current: 42.3, Duty: 91, Temperature: 0.2}},
			{buf: VehicleState{XDistance: 2.45, Current: 4.3, Duty: 255, Temperature: -10.2}, result: VehicleState{XDistance: 2.45, Current: 4.3, Duty: 255, Temperature: -10.2}},
		}

		for _, testCase := range cases {
			got := TestGetVehicleState(testCase.buf)

			if got != testCase.result {
				t.Fatalf("Wanted %f, got %f", testCase.result.Current, got.Current)
			}
		}

	})
}

func TestConversionBytes(t *testing.T) {
	t.Run("struct conversion", func(t *testing.T) {
		cases := []ConversionCaseBytes{
			{buf: []byte{154, 153, 153, 153, 153, 153, 3, 64, 51, 51, 51, 51, 51, 51, 17, 64, 1, 102, 102, 102, 102, 102, 102, 36, 64}, result: VehicleState{XDistance: 2.45, Current: 4.3, Duty: 1, Temperature: 10.2}},
		}

		for _, testCase := range cases {
			got := GetVehicleState(testCase.buf)

			if got != testCase.result {
				t.Fatalf("Wanted %f, got %f", testCase.result.Current, got.Current)
			}
		}

	})
}
