package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	// var dataVehicle []byte = []byte{154, 153, 153, 153, 153, 153, 3, 64, 51, 51, 51, 51, 51, 51, 17, 64, 1, 102, 102, 102, 102, 102, 102, 36, 64, 113, 61, 10, 215, 163, 240, 52, 64, 195, 245, 40, 92, 143, 194, 29, 64, 81, 41, 92, 143, 194, 245, 8, 77, 64, 154, 153, 153, 153, 153, 153, 3, 64, 51, 51, 51, 51, 51, 51, 17, 64, 1, 102, 102, 102, 102, 102, 102, 36, 64, 11, 11, 11}

	// result, err := Decode1(dataVehicle)
	// if err != nil {
	// 	fmt.Println("Error")
	// }

	// switch obj := result.(type) {
	// case []VehicleState:
	// 	fmt.Println("VehicleState", obj)
	// default:
	// 	fmt.Println("Don't know", obj)
	// }

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error al cargar archivo .env")
	}

	hilHandler := NewHilHandler()
	http.HandleFunc(os.Getenv("PATH"), func(w http.ResponseWriter, r *http.Request) { handle(w, r, hilHandler) })

	fmt.Println("Listening in", os.Getenv("SERVER_ADDR")+os.Getenv("PATH"))

	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_ADDR"), nil))

}

func handle(w http.ResponseWriter, r *http.Request, hilHandler *HilHandler) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	remoteHost, _, errSplit := net.SplitHostPort(r.RemoteAddr)
	if errSplit != nil {
		log.Println("Error spliting IP:", errSplit)
		return
	}

	if hilHandler.frontConn == nil && remoteHost == "127.0.0.1" { //TODO: Establish IP
		hilHandler.SetFrontConn(conn)
		fmt.Println(hilHandler, conn.RemoteAddr())
	}

	if hilHandler.hilConn == nil && remoteHost == "127.0.0.2" { //TODO: Establish IP from Hil
		hilHandler.SetHilConn(conn)
		fmt.Println(hilHandler, conn.RemoteAddr())
	}

	if hilHandler.frontConn != nil && hilHandler.hilConn != nil {
		errReady := hilHandler.frontConn.WriteMessage(websocket.TextMessage, []byte("Back-end is ready!"))
		if errReady != nil {
			log.Println("Error sending ready message:", errReady)
			return
		}
		hilHandler.StartIDLE()
	}
}
