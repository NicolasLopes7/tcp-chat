package protocol

import (
	"encoding/json"
	"fmt"
)

type MessageType uint8

const (
	SetName = iota
	SendMessage
	Logout
	Ping
)

type Message struct {
	Type    MessageType
	Payload string
}

func NewMessage(t MessageType, p string) *Message {
	return &Message{Type: t, Payload: p}
}

func (m *Message) ToBytes() []byte {
	b, _ := json.Marshal(m)

	return b
}

func (m *Message) ToString() string {
	return fmt.Sprintf("Type: %d\nPayload: \"%s\"", m.Type, m.Payload)
}

func ParseMessage(b []byte) *Message {
	var m *Message

	json.Unmarshal(b, &m)

	return m
}
