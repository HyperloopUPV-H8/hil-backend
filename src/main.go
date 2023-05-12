package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"main/utilities"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error al cargar archivo .env")
	}

	num := utilities.ConvertFloat64ToBytes(3.14)
	fmt.Println(num)
	fmt.Println(utilities.ConvertBytesToFloat64(num))

	mockBytes := utilities.CreateMockBytes()
	utilities.GetVehicleState(mockBytes[:])

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

	vehicleState := utilities.CreateVehicleState()

	fmt.Println(vehicleState)

	sendingVehicleStateJSON(conn)

}

func testVehicleStateToJSON() {
	vehicleState := utilities.CreateVehicleState()

	fmt.Println(vehicleState)

	vehicleStateJson, _ := json.Marshal(vehicleState)

	fmt.Println(vehicleStateJson)
	fmt.Println(string(vehicleStateJson))

	vehicleStateUnmarshalled := &utilities.VehicleState{}
	json.Unmarshal(vehicleStateJson, vehicleStateUnmarshalled)

	fmt.Println(vehicleStateUnmarshalled)
}

func sendingVehicleStateJSON(conn *websocket.Conn) {
	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				vehicleState := utilities.CreateVehicleState()
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

func receivingPerturbationData(conn *websocket.Conn) {
	go func() {
		for {
			perturbationData := &utilities.Perturbation{}
			conn.ReadJSON(perturbationData)
		}
	}()
}
