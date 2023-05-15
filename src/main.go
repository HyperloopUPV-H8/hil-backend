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

	// mockBytes := utilities.CreateMockBytes()
	// utilities.GetVehicleState(mockBytes[:])

	http.HandleFunc(os.Getenv("PATH"), handleWebSocket)

	fmt.Println("Listening in", os.Getenv("SERVER_ADDR"))

	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_ADDR"), nil))

}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool {
		//origin := r.Header.Get("Origin")
		//return origin == "http://127.0.0.1:5173/" || origin == "http://10.236.42.103:5173/"
		return true
	} //TODO: Check it the origin is correct

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}

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
	go func() {
		for range ticker.C {
			vehicleState := utilities.CreateVehicleState()

			errMarshal := conn.WriteJSON(vehicleState)

			if errMarshal != nil {
				log.Println("Error marshalling:", errMarshal)
				return
			}

			fmt.Println("struct sent!", vehicleState)
		}
	}()
}

func receivingPerturbationData(conn *websocket.Conn) {
	for {
		perturbationData := &utilities.PerturbationArray{}
		conn.ReadJSON(perturbationData)
	}
}

func createCloseHandler() (func(), <-chan bool) {
	done := make(chan bool)

	return func() {
		done <- true
	}, done
}
