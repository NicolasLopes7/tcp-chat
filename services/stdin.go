package services

import (
	"fmt"
	"strings"

	"github.com/NicolasLopes7/tcp-chat/network"
	"github.com/NicolasLopes7/tcp-chat/protocol"
)

type AbstractStdinConsumer interface {
	Parse(inputChan chan string) error
}

type StdinConsumer struct{}

func (sc *StdinConsumer) Parse(inputChan chan string, writer network.AbstractWriter) error {
	for {
		select {
		case line := <-inputChan:
			err := writer.Write(GetCommandOrMessage(line))
			if err != nil {
				fmt.Println("Error: ", err)
				return err
			}
		}
	}
}

func GetCommandOrMessage(line string) *protocol.Message {
	parts := strings.Split(line, " ")
	command := parts[0]

	switch command {
	case "/list":
		return protocol.NewMessage(protocol.ListUsers, "")
	case "/kick":
		return protocol.NewMessage(protocol.KickUser, parts[1])
	default:
		return protocol.NewMessage(protocol.SendMessage, line)
	}
}
