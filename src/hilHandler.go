package main

import (
	"encoding/json"
	"fmt"
	"log"
	"main/conversions"
	"main/models"
	"main/mvp1/utilities"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type HilHandler struct {
	frontConn *websocket.Conn
	hilConn   *websocket.Conn
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
		switch msg {
		case "start_simulation":
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
	dataChan := make(chan models.VehicleState)
	orderChan := make(chan models.Order)
	fmt.Println("Simulation state")

	hilHandler.startListeningData(dataChan, errChan, done)
	hilHandler.startSendingData(dataChan, errChan, done)

	hilHandler.startListeningOrders(orderChan, errChan, done)
	hilHandler.startSendingOrders(orderChan, errChan, done)

	err := <-errChan
	done <- struct{}{}

	return err
}

func (hilHandler *HilHandler) startSendingData(dataChan <-chan models.VehicleState, errChan <-chan error, done <-chan struct{}) {
	go func() {

		for {
			select {
			case <-done:
				return
			case <-errChan:
				return
			case data := <-dataChan:
				errMarshal := hilHandler.frontConn.WriteJSON(data)
				if errMarshal != nil {
					log.Println("Error marshalling:", errMarshal)
					return
				}
			}

		}
	}()
}

func (hilHandler *HilHandler) mockingSendVehicleState() {
	ticker := time.NewTicker(2 * time.Second)
	for range ticker.C {
		vehicleState := utilities.RandomVehicleState()

		errMarshal := hilHandler.frontConn.WriteJSON(vehicleState)

		if errMarshal != nil {
			log.Println("Error marshalling:", errMarshal)
			return
		}
	}
}

func (hilHandler *HilHandler) startListeningOrders(orderChan chan<- models.Order, errChan chan<- error, done <-chan struct{}) {
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				_, msg, errReadJSON := hilHandler.frontConn.ReadMessage()
				stringMsg := string(msg)
				if errReadJSON != nil {
					fmt.Println("err: ", errReadJSON)
					break
				}
				if strings.HasPrefix(stringMsg, "{\"id\":") {
					var order models.ControlOrder = models.ControlOrder{}
					errJSON := json.Unmarshal(msg, &order)
					if errJSON != nil {
						fmt.Println("err: ", errJSON)
						break
					}
					orderChan <- order
				} else if strings.HasPrefix(stringMsg, "[{\"id\":") {
					var orders models.FormData = models.FormData{}
					errJSON := json.Unmarshal(msg, &orders)
					if errJSON != nil {
						fmt.Println("err: ", errReadJSON)
						break
					}
					frontOrders := conversions.ConvertFormDataToOrders(orders)
					for _, frontOrder := range frontOrders {
						orderChan <- frontOrder
					}
				} else {
					fmt.Println("It is not an order")
				}
			}

		}
	}()
}

func (hilHandler *HilHandler) startListeningData(dataChan chan<- models.VehicleState, errChan chan<- error, done <-chan struct{}) {
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				_, msg, err := hilHandler.hilConn.ReadMessage()
				if err != nil {
					errChan <- err
					break
				}
				data, errDecoding := Decode(msg)
				if errDecoding != nil {
					fmt.Println("Error decoding: ", errDecoding)
					continue
				}

				switch decodedData := data.(type) {
				case []models.VehicleState:
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

func (hilHandler *HilHandler) startSendingOrders(orderChan <-chan models.Order, errChan <-chan error, done <-chan struct{}) {
	go func() {

		for {
			select {
			case <-done:
				return
			case <-errChan:
				return
			case order := <-orderChan:
				var orderArray []models.Order = []models.Order{order} //FIXME, from now it sends the order when it is received, to be defined if send several in same array
				encodedOrder := Encode(orderArray)
				errMarshal := hilHandler.hilConn.WriteMessage(websocket.BinaryMessage, encodedOrder)
				if errMarshal != nil {
					log.Println("Error marshalling:", errMarshal)
					return
				}
			}

		}
	}()
}
