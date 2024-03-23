package protocol

import (
	"encoding/json"
	"fmt"
	"net"
)

type MessageType uint8

const (
	SetName = iota
	SendMessage
	Logout
	Ping
	Die
	ListUsers
	KickUser
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

func ReadMessage(conn *net.Conn) (*Message, error) {
	buffer := make([]byte, 1024)
	n, err := (*conn).Read(buffer)
	if err != nil {
		return nil, err
	}

	message := ParseMessage(buffer[:n])
	return message, nil
}
