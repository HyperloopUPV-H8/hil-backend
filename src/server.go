package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	handleConn func(conn *websocket.Conn)
}

func NewServer() Server {
	return Server{
		handleConn: func(conn *websocket.Conn) {},
	}

}

func (server *Server) SetConnHandler(handler func(conn *websocket.Conn)) { //TODO: add as prop handler
	server.handleConn = handler
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	server.handleConn(conn)
}
