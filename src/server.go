package main

import (
	"net/http"

	"github.com/gorilla/websocket"
	trace "github.com/rs/zerolog/log"
)

type Server struct {
	handleConn func(conn *websocket.Conn)
}

func NewServer() Server {
	return Server{
		handleConn: func(conn *websocket.Conn) {},
	}

}

func (server *Server) SetConnHandler(handler func(conn *websocket.Conn)) {
	server.handleConn = handler
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		trace.Error().Err(err).Msg("Error upgrading connection")
		return
	}
	server.handleConn(conn)
}
