package main

import (
	"fmt"
	"log"
	"main/mvp1/utilities"
	"time"

	"github.com/gorilla/websocket"
)

func NewHilHandler() *HilHandler {
	return &HilHandler{}
}

func (hilHandler *HilHandler) SetFrontConn(conn *websocket.Conn) {
	hilHandler.frontConn = conn
}

func (hilHandler *HilHandler) SetHilConn(conn *websocket.Conn) {
	hilHandler.hilConn = conn
}

func (hilHandler *HilHandler) StartIDLE() {
	for {
		_, msgByte, err := hilHandler.frontConn.ReadMessage()
		if err != nil {
			log.Fatalf("error receiving message in IDLE: ", err)
		}
		msg := string(msgByte)
		fmt.Println(msg)
		switch msg {
		case "start_simulation":
			err := hilHandler.startSimulationState()

			if err != nil {
				return
			}
		}

		//fmt.Println("Llega a fin del IDLE") No llega porque readMessage deja bloqueado

	}
}

func (hilHandler *HilHandler) startSimulationState() error {
	errChan := make(chan error)
	done := make(chan struct{})
	dataChan := make(chan VehicleState)
	//orderChan := make(chan Order)
	// Recibo info del HIL y la envio al front
	go hilHandler.startSendingData(dataChan, errChan, done)

	// Recibo ordenes del front y la envio al HIL
	//go hilHandler.startTransmittingOrders(orderChan, errChan, done)

	err := <-errChan
	//Aquí se bloquea si llega error?
	done <- struct{}{}

	return err
}

func (hilHandler *HilHandler) startSendingData(dataChan <-chan VehicleState, errChan <-chan error, done <-chan struct{}) {
	ticker := time.NewTicker(2 * time.Second)
	go func() {
		for {
			select {
			case <-done:
				return
			case data := <-dataChan: //TODO: Define msg origin, now it is a mock
				errMarshal := hilHandler.frontConn.WriteJSON(data)
				if errMarshal != nil {
					log.Println("Error marshalling:", errMarshal)
					return
				}

				fmt.Println("struct sent!", data)
			default:
				for range ticker.C {
					//FIXME: Only for mocking
					vehicleState := utilities.RandomVehicleState()

					errMarshal := hilHandler.frontConn.WriteJSON(vehicleState)

					if errMarshal != nil {
						log.Println("Error marshalling:", errMarshal)
						return
					}

					fmt.Println("struct sent!", vehicleState)
				}
			}

		}
	}()
}

// func (hilHandler *HilHandler) startTransmittingOrders(dataChan <-chan VehicleState, errChan <-chan error, done <-chan struct{}) {
// }

// func (hilHandler *HilHandler) startTransmittingData(dataChan chan<- SimulationData, errChan chan<- error, done <-chan struct{}) {
//     for {
//         select {
//         case <-done:
//             break
//         default:
//             buf, err := hilHandler.hilConn.Read()
//             if err != nil {
//                 errChan<-err
//                 break
//             }

//             data, err := hilHandler.parser.Decode(buf)
//             if err != nil {
//                 continue
//             }

//             dataChan <- data
//         }

//     }
// }
