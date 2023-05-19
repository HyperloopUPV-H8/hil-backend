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
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error al cargar archivo .env")
	}

	http.HandleFunc(os.Getenv("PATH"), handle)

	fmt.Println("Listening in", os.Getenv("SERVER_ADDR"))

	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_ADDR"), nil))

}

func handle(w http.ResponseWriter, r *http.Request) {

	hilHandler := NewHilHandler()

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
	if remoteHost == "127.0.0.1" { //TODO: Establish IP
		hilHandler.SetFrontConn(conn)
	}

	if remoteHost == "127.0.0.1" { //TODO: Establish IP from Hil
		hilHandler.SetHilConn(conn)
		fmt.Println(hilHandler, r.RemoteAddr, conn.RemoteAddr())
	}

	if hilHandler.frontConn != nil && hilHandler.hilConn != nil {
		//byteArray := []byte{97, 98, 99, 100, 101, 102}
		//var msg []byte = []byte[11111,11]
		//[]byte("Back-end is ready!")

		errReady := hilHandler.frontConn.WriteMessage(websocket.TextMessage, []byte("Back-end is ready!"))
		if errReady != nil {
			log.Println("Error sending ready message:", errReady)
			return
		}
		hilHandler.StartIDLE()
		//mvp1.SendingVehicleStateJSON(conn)
	}
}
