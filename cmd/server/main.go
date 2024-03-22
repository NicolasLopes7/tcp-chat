package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/NicolasLopes7/tcp-chat/protocol"
	"github.com/NicolasLopes7/tcp-chat/state"
)

type Server struct {
	listener net.Listener
	clients  *state.ClientStore
	mutex    sync.RWMutex
	shutdown bool
}

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	clientStore := state.NewClientStore()

	server := &Server{
		listener: listener,
		clients:  clientStore,
	}

	go func() {
		<-signalChan
		fmt.Println("\nShutdown signal received. Closing all connections...")

		server.mutex.Lock()
		for addr, client := range clientStore.Clients {
			fmt.Printf("Closing connection with %s\n", addr)
			(*client.Conn).Write(protocol.NewMessage(protocol.Die, "Server is shutting down").ToBytes())
			(*client.Conn).Close()
		}
		server.shutdown = true
		server.mutex.Unlock()

		server.listener.Close()
		fmt.Println("All connections closed. Shutting down the server...")
		os.Exit(0)
	}()

	fmt.Println("🚀 Server started. Listening connections on port 8080")

	for {
		conn, err := server.listener.Accept()
		if err != nil {
			if server.shutdown {
				return
			}
			fmt.Println("Error: ", err)
			return
		}

		server.clients.Add(conn.RemoteAddr().String(), &state.Client{Conn: &conn, Name: ""})
		go server.Handle(conn)
	}
}

func (s *Server) Handle(conn net.Conn) {
	defer func() {
		s.clients.Delete(conn.RemoteAddr().String())
		conn.Close()
	}()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil || n == 0 {
			if s.shutdown {
				return
			}
			s.clients.Delete(conn.RemoteAddr().String())
			return
		}

		m := protocol.ParseMessage(buffer[:n])

		switch m.Type {
		case protocol.SetName:
			fmt.Printf("%s connected as %s\n", conn.RemoteAddr().String(), m.Payload)
			s.clients.Add(conn.RemoteAddr().String(), &state.Client{Conn: &conn, Name: m.Payload})

		case protocol.SendMessage:
			if client, ok := s.clients.Get(conn.RemoteAddr().String()); ok {
				fmt.Println(client.Name + ": " + m.Payload)
			}

		case protocol.ListUsers:
			fmt.Println("List of users:")
			for addr, client := range s.clients.Clients {
				fmt.Printf("- %s (%s)\n", client.Name, addr)
			}

		case protocol.Logout:
			fmt.Printf("%s disconnected\n", conn.RemoteAddr().String())
			s.clients.Delete(conn.RemoteAddr().String())
			return

		case protocol.KickUser:
			for addr, client := range s.clients.Clients {
				if client.Name == m.Payload {
					fmt.Printf("Kicking %s\n", addr)
					(*client.Conn).Write(protocol.NewMessage(protocol.Die, "You were kicked").ToBytes())
					(*client.Conn).Close()
					s.clients.Delete(addr)
				}
			}
		}
	}
}
