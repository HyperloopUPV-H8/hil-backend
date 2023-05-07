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
	XDistance float64 //It goes between 22mm and 10 mm
	XRotation float64
	YRotation float64
	ZRotation float64
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

	http.HandleFunc(os.Getenv("PATH"), handleWebSocket1)

	// Iniciar el servidor HTTP en el puerto 8010
	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_ADDR"), nil))
}

func handleWebSocket1(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error al actualizar la conexión:", err)
		return
	}
	defer conn.Close()

	msg := []byte("message")

	err = conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		log.Println("Error al escribir el mensaje:", err)
	}

	fmt.Print("Hola")
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error al actualizar la conexión:", err)
		return
	}
	defer conn.Close()

	// Leer los mensajes enviados por el cliente
	//Pruebas lectura de mensaje
	// for {
	// 	_, message, err := conn.ReadMessage()
	// 	json.RawMessage
	// 	var myStruct Message

	// 	json.Unmarshal(message, myStruct)

	// 	if err != nil {
	// 		log.Println("Error al leer el mensaje:", err)
	// 		break
	// 	}
	// 	log.Printf("Mensaje recibido: %s\n", message)

	// 	// Responder al cliente con el mismo mensaje recibido
	// 	err = conn.WriteMessage(websocket.TextMessage, message)
	// 	if err != nil {
	// 		log.Println("Error al escribir el mensaje:", err)
	// 		break
	// 	}
	// }

}

func sendingVehicleState(conn *websocket.Conn, vehicleState VehicleState) {
	buf, err := json.Marshal(vehicleState)

	if err != nil {
		log.Println("Error marshalling the JSON:", err)
	}

	errWriting := conn.WriteJSON(buf)

	if errWriting != nil {
		log.Println("Error sending the JSON:", err)
	}
}

func sendingPerturbationData(conn *websocket.Conn) {
	ticker := time.NewTicker(time.Second * 2)

	go func() {
		for range ticker.C {
			perturbationData := createPerturbationData()
			buf, err := json.Marshal(perturbationData)

			if err != nil {
				log.Println("Error marshalling the JSON:", err)
			}

			errWriting := conn.WriteJSON(buf)

			if errWriting != nil {
				log.Println("Error sending the JSON:", err)
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
