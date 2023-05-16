package utilities

import (
	"log"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

func RandomPerturbationData() Perturbation {
	perturbationData := &Perturbation{}
	perturbationData.TypeP = SelectPerturbationType()
	perturbationData.Value = rand.Intn(10)

	return *perturbationData
}

func SelectPerturbationType() string {
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

func SendingPerturbationData(conn *websocket.Conn) {
	ticker := time.NewTicker(time.Second * 2)

	go func() {
		for range ticker.C {
			perturbationData := RandomPerturbationData()
			errWriting := conn.WriteJSON(perturbationData)

			if errWriting != nil {
				log.Println("Error sending the JSON:", errWriting)
			}
		}
	}()
}
