package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/pelletier/go-toml/v2"
	trace "github.com/rs/zerolog/log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	config := getConfig("./config.toml")

	hilHandler := NewHilHandler()

	http.HandleFunc(config.Path, func(w http.ResponseWriter, r *http.Request) { handle(w, r, hilHandler, config.Addresses) })
	fmt.Println("Listening in", config.Addresses.Server_addr+config.Path)
	log.Fatal(http.ListenAndServe(config.Addresses.Server_addr, nil))

}

func handle(w http.ResponseWriter, r *http.Request, hilHandler *HilHandler, addressesConfgi AddressesCongif) {
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
	if hilHandler.frontConn == nil && remoteHost == addressesConfgi.Frontend {
		hilHandler.SetFrontConn(conn)
		fmt.Println("Frontened connected: ", hilHandler, conn.RemoteAddr())
	}

	if hilHandler.hilConn == nil && remoteHost == addressesConfgi.Hil {
		hilHandler.SetHilConn(conn)
		fmt.Println("HIL connected: ", hilHandler, conn.RemoteAddr())
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

func getConfig(path string) Config {
	configFile, fileErr := os.ReadFile(path)

	if fileErr != nil {
		trace.Fatal().Stack().Err(fileErr).Msg("error reading config file")
	}
	var config Config
	tomlErr := toml.Unmarshal(configFile, &config)
	if tomlErr != nil {
		trace.Fatal().Stack().Err(tomlErr).Msg("error unmarshalling toml")
	}
	return config
}
