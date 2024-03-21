package main

import (
	"fmt"
	"net"

	"github.com/NicolasLopes7/tcp-chat/protocol"
	"github.com/NicolasLopes7/tcp-chat/state"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	cs := state.NewClientStore()
	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		go handle(conn, cs)
	}
}

func handle(conn net.Conn, cs *state.ClientStore) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil || n == 0 {
			fmt.Println("Error: ", err)
			cs.Delete(conn.RemoteAddr().String())
			return
		}

		m := protocol.ParseMessage(buffer[:n])

		switch m.Type {
		case protocol.SetName:
			fmt.Printf("%s connected as %s\n", conn.RemoteAddr().String(), m.Payload)
			cs.Add(conn.RemoteAddr().String(), m.Payload)

		case protocol.SendMessage:
			if name, ok := cs.Get(conn.RemoteAddr().String()); ok {
				fmt.Println(name + ": " + m.Payload)
			}

		case protocol.Logout:
			fmt.Printf("%s disconnected\n", conn.RemoteAddr().String())
			cs.Delete(conn.RemoteAddr().String())
			return
		}
	}
}
