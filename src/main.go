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

	http.HandleFunc(os.Getenv("PATH"), handleWebSocket)

	fmt.Println("Listening in", os.Getenv("SERVER_ADDR"))

	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_ADDR"), nil))

}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool {
		//TODO: Check it the origin is correct
		//origin := r.Header.Get("Origin")
		//return origin == "http://127.0.0.1:5173/" || origin == "http://10.236.42.103:5173/"
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}

	sendingVehicleStateJSON(conn)
	receivingStringMessage(conn)

}

func sendingVehicleStateJSON(conn *websocket.Conn) {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			//TODO: Define msg origin, now it is a mock
			vehicleState := utilities.RandomVehicleState()

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
	go func() {
		fmt.Println("Init receiving")
		for {
			//TODO: Define struct which it's gonna be received
			perturbationData := &utilities.PerturbationArray{}
			err := conn.ReadJSON(perturbationData)
			if err != nil {
				fmt.Println("Failed")
			}
			//TODO: What to do with the formData received
		}
	}()
}

func receivingStringMessage(conn *websocket.Conn) {
	go func() {
		fmt.Println("Init receiving")
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Failed")
			}
			fmt.Println(string(msg[:]))
		}
	}()
}

func testVehicleStateToJSON() {
	vehicleState := utilities.RandomVehicleState()

	fmt.Println(vehicleState)

	vehicleStateJson, _ := json.Marshal(vehicleState)

	fmt.Println(vehicleStateJson)
	fmt.Println(string(vehicleStateJson))

	vehicleStateUnmarshalled := &utilities.VehicleState{}
	json.Unmarshal(vehicleStateJson, vehicleStateUnmarshalled)

	fmt.Println(vehicleStateUnmarshalled)
}

func createCloseHandler() (func(), <-chan bool) {
	done := make(chan bool)

	return func() {
		done <- true
	}, done
}
