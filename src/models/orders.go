package models

import (
	"bytes"
	"encoding/binary"
)

type Order interface {
	Bytes() []byte
}

type FrontOrder struct {
	Kind    string  `json:"kind"`
	Payload float64 `json:"payload"`
}

func (order FrontOrder) Bytes() []byte {
	buf1 := []byte(order.Kind)
	var buf2 [8]byte
	binary.LittleEndian.PutUint64(buf2[:], uint64(order.Payload))
	return append(buf1, buf2[:]...)
}

func (order *FrontOrder) Read(data []byte) {
	reader := bytes.NewReader(data)
	binary.Read(reader, binary.LittleEndian, order)
}

type ControlOrder struct {
	Id    uint8 `json:"id"`
	State bool  `json:"state"`
}

func (order ControlOrder) Bytes() []byte {
	buf1 := []byte{order.Id}
	var booleanValue uint8
	if order.State {
		booleanValue = 1
	} else {
		booleanValue = 0
	}
	return append(buf1, booleanValue)
}

func (order *ControlOrder) Read(data []byte) {
	reader := bytes.NewReader(data)
	binary.Read(reader, binary.LittleEndian, order)
}
