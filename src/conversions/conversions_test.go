package conversions

import (
	"main/models"
	"reflect"
	"testing"
)

type ConversionCaseBytes struct {
	buf    []byte
	result []models.VehicleState
}

func TestConversionToVehicleStateArray(t *testing.T) {
	t.Run("struct conversion", func(t *testing.T) {
		cases := []ConversionCaseBytes{
			{buf: []byte{154, 153, 153, 153, 153, 153, 3, 64, 51, 51, 51, 51, 51, 51, 17, 64, 1, 102, 102, 102, 102, 102, 102, 36, 64, 113, 61, 10, 215, 163, 240, 52, 64, 195, 245, 40, 92, 143, 194, 29, 64, 81, 41, 92, 143, 194, 245, 8, 77, 64, 154, 153, 153, 153, 153, 153, 3, 64, 51, 51, 51, 51, 51, 51, 17, 64, 1, 102, 102, 102, 102, 102, 102, 36, 64, 11, 11, 11},
				result: []models.VehicleState{
					{YDistance: 2.45, Current: 4.3, Duty: 1, Temperature: 10.2},
					{YDistance: 20.94, Current: 7.44, Duty: 81, Temperature: 58.07},
					{YDistance: 2.45, Current: 4.3, Duty: 1, Temperature: 10.2}}},
		}

		for _, testCase := range cases {
			got, _ := GetAllVehicleStates(testCase.buf)

			if !reflect.DeepEqual(got, testCase.result) {
				t.Fatalf("Wanted %f, got %f", testCase.result[1].Current, got[1].Current)
			}
		}

	})
}
