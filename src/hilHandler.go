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

	//parser HilParser // Tiene Encode y Decode
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
	fmt.Println("IDLE")
	for {
		_, msgByte, err := hilHandler.frontConn.ReadMessage()
		if err != nil {
			log.Fatalf("error receiving message in IDLE: %s", err)
		}
		msg := string(msgByte)
		fmt.Println(msg)
		switch msg {
		case "start_simulation": //FIXME: Check in front when it sends the msg
			err := hilHandler.startSimulationState()

			if err != nil {
				return
			}
		}
	}
}

func (hilHandler *HilHandler) startSimulationState() error {
	errChan := make(chan error)
	done := make(chan struct{})
	dataChan := make(chan VehicleState)
	orderChan := make(chan Order)
	fmt.Println("Simulation state")

	//FIXME: IT WORKS PROPERLY
	hilHandler.startListeningData(dataChan, errChan, done)
	hilHandler.startSendingData(dataChan, errChan, done)

	// From front to HIL
	hilHandler.startListeningOrders(orderChan, errChan, done)
	hilHandler.startSendingOrders(orderChan, errChan, done)

	//FIXME: Waiting for error, and sending done to close the rest
	err := <-errChan
	done <- struct{}{}

	return err
}

func (hilHandler *HilHandler) startSendingData(dataChan <-chan VehicleState, errChan <-chan error, done <-chan struct{}) {
	go func() {

		for {
			select {
			case <-done:
				return
			case <-errChan: //FIXME
				return
			case data := <-dataChan:
				fmt.Println(data.Current)
				errMarshal := hilHandler.frontConn.WriteJSON(data)
				if errMarshal != nil {
					log.Println("Error marshalling:", errMarshal)
					return
				}

				fmt.Println("struct sent!", data)
			default:
			}

		}
	}()
}

// Only for mocking and tests for func StartSendingData
func (hilHandler *HilHandler) mockingSendVehicleState() {
	ticker := time.NewTicker(2 * time.Second)
	for range ticker.C {
		vehicleState := utilities.RandomVehicleState()

		errMarshal := hilHandler.frontConn.WriteJSON(vehicleState)

		if errMarshal != nil {
			log.Println("Error marshalling:", errMarshal)
			return
		}

		fmt.Println("struct sent!", vehicleState)
	}
}

func (hilHandler *HilHandler) startListeningOrders(orderChan chan<- Order, errChan chan<- error, done <-chan struct{}) {
	go func() {
		for {
			select {
			case <-done:
				return
				//case <-errChan: //FIXME
				//return

			default:
				var order Order
				errReadJSON := hilHandler.frontConn.ReadJSON(order)
				if errReadJSON != nil {
					errChan <- errReadJSON
					break
				}
				orderChan <- order

			}

		}
	}()
}

func (hilHandler *HilHandler) startListeningData(dataChan chan<- VehicleState, errChan chan<- error, done <-chan struct{}) {
	go func() {
		for {
			select {
			case <-done:
				return
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

				switch decodedData := data.(type) {
				case []VehicleState:
					for _, d := range decodedData {
						dataChan <- d
					}
				default:
					fmt.Println("Does NOT match any type (startListeningData): ", decodedData)
				}

			}

		}
	}()
}

func (hilHandler *HilHandler) startSendingOrders(orderChan <-chan Order, errChan <-chan error, done <-chan struct{}) {
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		for {
			select {
			case <-done:
				return
			case <-errChan: //FIXME
				return
			case order := <-orderChan:
				var orderArray []Order = []Order{order} //TODO, it is gonna use arrays or only a order
				encodedOrder := Encode(orderArray)
				errMarshal := hilHandler.hilConn.WriteMessage(websocket.BinaryMessage, encodedOrder)
				if errMarshal != nil {
					log.Println("Error marshalling:", errMarshal)
					return
				}
			default:
				for range ticker.C {
					//TODO: Send orders to HIL
				}
			}

		}
	}()
}
