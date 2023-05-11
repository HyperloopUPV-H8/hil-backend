package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

type Perturbation struct {
	TypeP string
	Value int
}

type VehicleState struct {
	XDistance float64 `json:"xDistance"` //It goes between 22mm and 10 mm
	XRotation float64 `json:"xRotation"`
	// YRotation float64
	// ZRotation float64
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error al cargar archivo .env")
	}

	vehicleState := createVehicleState()

	fmt.Println(vehicleState)

	vehicleStateJson, _ := json.Marshal(vehicleState)

	fmt.Println(vehicleStateJson)
	fmt.Println(string(vehicleStateJson))

	vehicleStateUnmarshalled := &VehicleState{}
	json.Unmarshal(vehicleStateJson, vehicleStateUnmarshalled)

	fmt.Println(vehicleStateUnmarshalled)

	http.HandleFunc(os.Getenv("PATH"), handleWebSocket)

	fmt.Println("Listening in", os.Getenv("SERVER_ADDR"))
	// Iniciar el servidor HTTP en el puerto 8010
	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_ADDR"), nil))

}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	vehicleState := createVehicleState()

	fmt.Println(vehicleState)

	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				vehicleState := createVehicleState()
				errMarshal := conn.WriteJSON(vehicleState)

				if errMarshal != nil {
					log.Println("Error marshalling:", errMarshal)
					return
				}

				fmt.Println("struct sent!", vehicleState)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

}

func createVehicleState() VehicleState { //TODO
	VehicleState := &VehicleState{}
	VehicleState.XDistance = float64(rand.Intn(12) + 10.0)
	VehicleState.XRotation = 0
	return *VehicleState
}

func sendingPerturbationData(conn *websocket.Conn) {
	ticker := time.NewTicker(time.Second * 2)

	go func() {
		for range ticker.C {
			perturbationData := createPerturbationData()
			errWriting := conn.WriteJSON(perturbationData)

			if errWriting != nil {
				log.Println("Error sending the JSON:", errWriting)
			}
		}
	}()
}

func createPerturbationData() Perturbation {
	perturbationData := &Perturbation{}
	perturbationData.TypeP = selectPerturbationType()
	perturbationData.Value = rand.Intn(10)

	return *perturbationData
}

func selectPerturbationType() string {
	perturbationType := rand.Intn(3)
	switch {
	case perturbationType == 0:
		return "x-axis"
	case perturbationType == 1:
		return "y-axis"
	case perturbationType == 2:
		return "z-axis"
	default:
		return ""
	}
}
