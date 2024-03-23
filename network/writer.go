package network

import (
	"fmt"
	"net"

	"github.com/NicolasLopes7/tcp-chat/protocol"
)

type AbstractWriter interface {
	Write(message *protocol.Message) error
}

type WriterService struct {
	conn *net.Conn
}

func (ws *WriterService) Write(message *protocol.Message) error {

	_, err := (*ws.conn).Write(message.ToBytes())

	if err != nil {
		fmt.Println("Error: ", err)
	}

	return nil
}
