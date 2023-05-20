package main

import (
	"fmt"
	"log"
	"main/mvp1/utilities"
	"time"

	"github.com/gorilla/websocket"
)

type HilHandler struct {
	frontConn *websocket.Conn
	hilConn   *websocket.Conn

	parser HilParser // Tiene Encode y Decode
}

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
			log.Fatalf("error receiving message in IDLE: %s", err)
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
	orderChan := make(chan Order)
	// Recibo info del HIL y la envio al front
	go hilHandler.startSendingData(dataChan, errChan, done)
	//FIXME: Is it necesary to put go an inside the go func?

	// Recibo ordenes del front y la envio al HIL
	go hilHandler.startTransmittingOrders(orderChan, errChan, done)

	//FIXME: Is it waiting for both? Or only for error?
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
			case <-errChan: //FIXME
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

func (hilHandler *HilHandler) startTransmittingOrders(orderChan <-chan Order, errChan <-chan error, done <-chan struct{}) {
	go func() {
		for {
			select {
			case <-done:
				return
			case <-errChan: //FIXME
				return
			// case order := <-orderChan: //TODO: Define msg origin, now it is a mock
			// 	//Encode
			// 	errMarshal := hilHandler.hilConn.WriteMessage(websocket.BinaryMessage, order)
			// 	if errMarshal != nil {
			// 		log.Println("Error marshalling:", errMarshal)
			// 		return
			// 	}

			default:

			}

		}
	}()
}

func (hilHandler *HilHandler) startTransmittingData(dataChan chan<- VehicleState, errChan chan<- error, done <-chan struct{}) {
	for {
		select {
		case <-done:
			break
		default:
			_, buf, err := hilHandler.hilConn.ReadMessage()
			if err != nil {
				errChan <- err
				break
			}
			data, errDecoding := Decode(buf) //TODO: Decode
			if errDecoding != nil {
				fmt.Println("Error decoding: ", errDecoding)
				continue
			}
			var decodedData []VehicleState = data.([]VehicleState) //FIXME: Change depends on type

			for _, d := range decodedData {
				dataChan <- d
			}

		}

	}
}
