package main

import (
	"encoding/json"
	"fmt"
	"log"
	"main/conversions"
	"main/models"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	trace "github.com/rs/zerolog/log"
)

const STOP_MSG = "finish_simulation"

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
	trace.Info().Msg("IDLE")
	for {
		_, msgByte, err := hilHandler.frontConn.ReadMessage()
		if err != nil {
			log.Fatalf("error receiving message in IDLE: %s", err)
		}
		msg := string(msgByte)
		switch msg {
		case "start_simulation":
			err := hilHandler.startSimulationState()
			trace.Info().Msg("IDLE")

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
	stopChan := make(chan struct{})
	trace.Info().Msg("Simulation state")

	hilHandler.startListeningData(dataChan, errChan, done)
	hilHandler.startSendingData(dataChan, errChan, done)

	hilHandler.startListeningOrders(orderChan, errChan, done, stopChan)
	hilHandler.startSendingOrders(orderChan, errChan, done)

	for {
		select {
		case err := <-errChan:
			close(done)
			return err
		case <-stopChan:
			close(done)
			return nil
		default:
		}
	}
}

func (hilHandler *HilHandler) startSendingData(dataChan <-chan models.VehicleState, errChan chan<- error, done <-chan struct{}) {
	go func() {

		for {
			select {
			case <-done:
				return
			case data := <-dataChan:
				errMarshal := hilHandler.frontConn.WriteJSON(data)
				if errMarshal != nil {
					trace.Error().Err(errMarshal).Msg("Error marshalling")
					errChan <- errMarshal
					return
				}
			default:
			}

		}
	}()
}

func (hilHandler *HilHandler) mockingSendVehicleState() {
	ticker := time.NewTicker(2 * time.Second)
	for range ticker.C {
		vehicleState := models.RandomVehicleState()

		errMarshal := hilHandler.frontConn.WriteJSON(vehicleState)

		if errMarshal != nil {
			log.Println("Error marshalling:", errMarshal)
			return
		}
	}
}

func (hilHandler *HilHandler) startListeningOrders(orderChan chan<- models.Order, errChan chan<- error, done <-chan struct{}, stopChan chan<- struct{}) {
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				_, msg, errReadJSON := hilHandler.frontConn.ReadMessage() //FIXME: This block and can't return the done at the moment
				stringMsg := string(msg)
				if errReadJSON != nil {
					trace.Error().Err(errReadJSON).Msg("Error reading message from frontend")
					errChan <- errReadJSON
					break
				}
				if stringMsg == STOP_MSG {
					trace.Info().Msg("Stop simulation")
					stopChan <- struct{}{}
					return
				} else if !addOrderToChan(msg, orderChan) {
					trace.Warn().Msg("It is not an order")
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
				_, msg, err := hilHandler.hilConn.ReadMessage() //FIXME, it get block when done is close, if not new msg, it get stuck
				if err != nil {
					errChan <- err
					break
				}
				data, errDecoding := Decode(msg)
				if errDecoding != nil {
					trace.Error().Err(errDecoding).Msg(fmt.Sprintf("Error decoding: %v", errDecoding))
					continue
				}

				switch decodedData := data.(type) {
				case []models.VehicleState:
					for _, d := range decodedData {
						dataChan <- d
					}
				default:
					trace.Warn().Msg(fmt.Sprintf("Does NOT match any type (startListeningData): %v", decodedData))
				}

			}

		}
	}()
}

func (hilHandler *HilHandler) startSendingOrders(orderChan <-chan models.Order, errChan chan<- error, done <-chan struct{}) {
	go func() {

		for {
			select {
			case <-done:
				return
			case order := <-orderChan:
				encodedOrder := order.Bytes()
				errMarshal := hilHandler.hilConn.WriteMessage(websocket.BinaryMessage, encodedOrder)
				if errMarshal != nil {
					log.Println("Error marshalling:", errMarshal)
					errChan <- errMarshal
					return
				}
			default:
			}

		}
	}()
}

func prepareFormOrder(msg []byte, orderChan chan<- models.Order) error {
	var orders models.FormData = models.FormData{}
	errJSON := json.Unmarshal(msg, &orders)
	if errJSON != nil {
		trace.Error().Err(errJSON).Msg("Error unmarshalling Form Data")
		return errJSON
	}
	formOrders := conversions.ConvertFormDataToOrders(orders)
	for _, formOrder := range formOrders {
		orderChan <- formOrder
	}
	return nil
}

func addOrderToChan(msg []byte, orderChan chan<- models.Order) bool {
	stringMsg := string(msg)
	if strings.HasPrefix(stringMsg, "{\"id\":") {
		var order models.ControlOrder = models.ControlOrder{}
		errJSON := json.Unmarshal(msg, &order)
		if errJSON != nil {
			trace.Error().Err(errJSON).Msg("Error unmarshalling Control Order")
			return true
		}
		orderChan <- order
	} else if strings.HasPrefix(stringMsg, "[{\"id\":") {
		errJSON := prepareFormOrder(msg, orderChan)
		if errJSON != nil {
			trace.Error().Err(errJSON).Msg("Error unmarshalling Form Data")
			return true
		}
	} else {
		return false
	}
	return true

}
