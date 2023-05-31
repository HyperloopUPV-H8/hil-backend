package main

import (
	"flag"
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

var traceLevel = flag.String("trace", "info", "set the trace level (\"fatal\", \"error\", \"warn\", \"info\", \"debug\", \"trace\")")
var traceFile = flag.String("log", "trace.json", "set the trace log file")

func main() {
	traceFile := initTrace(*traceLevel, *traceFile)
	defer traceFile.Close()

	config := getConfig("./config.toml")

	hilHandler := NewHilHandler()

	server := NewServer()

	server.SetConnHandler(func(conn *websocket.Conn) {
		remoteHost, _, errSplit := net.SplitHostPort(conn.RemoteAddr().String())
		if errSplit != nil {
			log.Println("Error spliting IP:", errSplit)
			return
		}

		if hilHandler.frontConn == nil && remoteHost == config.Addresses.Frontend {
			hilHandler.SetFrontConn(conn)
			frontMsg := fmt.Sprintf("Frontened connected: %v %v", hilHandler, conn.RemoteAddr())
			trace.Info().Msg(frontMsg)

		}

		if hilHandler.hilConn == nil && remoteHost == config.Addresses.Hil {
			hilHandler.SetHilConn(conn)
			hilMsg := fmt.Sprintf("HIL connected: %v %v", hilHandler, conn.RemoteAddr())
			trace.Info().Msg(hilMsg)
		}

		if hilHandler.frontConn != nil && hilHandler.hilConn != nil {
			errReady := hilHandler.frontConn.WriteMessage(websocket.TextMessage, []byte("Back-end is ready!"))
			if errReady != nil {
				log.Println("Error sending ready message:", errReady)
				return
			}
			hilHandler.StartIDLE()
		}
	})

	http.Handle(config.Path, &server)

	trace.Info().Msg("Listening in " + config.Addresses.Server_addr + config.Path)
	log.Fatal(http.ListenAndServe(config.Addresses.Server_addr, nil))

}

func handle(w http.ResponseWriter, r *http.Request, hilHandler *HilHandler, addressesConfig AddressesConfig) {
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

	if hilHandler.frontConn == nil && remoteHost == addressesConfig.Frontend {
		hilHandler.SetFrontConn(conn)
		frontMsg := fmt.Sprintf("Frontened connected: %v %v", hilHandler, conn.RemoteAddr())
		trace.Info().Msg(frontMsg)

	}

	if hilHandler.hilConn == nil && remoteHost == addressesConfig.Hil {
		hilHandler.SetHilConn(conn)
		hilMsg := fmt.Sprintf("HIL connected: %v %v", hilHandler, conn.RemoteAddr())
		trace.Info().Msg(hilMsg)
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
