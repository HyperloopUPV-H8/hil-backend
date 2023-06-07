package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"main/conversions"
	"main/models"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	trace "github.com/rs/zerolog/log"
)

const START_MSG = "start_simulation"
const FINISH_SIMULATION = "finish_simulation"

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
			trace.Error().Err(err).Msg("error receiving message in IDLE")
		} else {
			msg := string(msgByte)
			switch msg {
			case START_MSG:

				errStartingHIL := hilHandler.informStartSimulation()
				if errStartingHIL != nil {
					trace.Error().Err(errStartingHIL).Msg("error informing HIL to start simulation")
					break
				}

				err := hilHandler.startSimulationState()
				trace.Info().Msg("IDLE")

				if err != nil {
					return
				}
			}
		}
	}
}

func (hilHandler *HilHandler) informStartSimulation() error {
	errStarting := hilHandler.hilConn.WriteMessage(websocket.BinaryMessage, []byte(START_MSG))
	if errStarting != nil {
		trace.Error().Err(errStarting).Msg("Error sending message of starting simultaion to HIL")
		return errStarting
	}
	_, msgByte, err := hilHandler.hilConn.ReadMessage()
	if err != nil {
		trace.Error().Err(err).Msg("error receiving message in IDLE")
		return err
	} else if string(msgByte) == START_MSG {
		return nil
	} else {
		errReceiving := errors.New("not received correct message from hil")
		trace.Error().Err(errReceiving).Msg("error receiving message in IDLE")
		return errReceiving
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
					return //FIXME: Before it was a break
				}
				if stringMsg == FINISH_SIMULATION {
					trace.Info().Msg("Finish simulation")
					errStoping := hilHandler.hilConn.WriteMessage(websocket.BinaryMessage, []byte(FINISH_SIMULATION))
					if errStoping != nil {
						trace.Error().Err(errStoping).Msg("Error sending finish simulation to HIL")
						errChan <- errStoping
					}
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
