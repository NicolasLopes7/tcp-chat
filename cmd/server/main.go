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
	cache    *state.UserCache
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

	userCache := state.NewUserCache()

	server := &Server{
		listener: listener,
		cache:    userCache,
	}

	go func() {
		<-signalChan
		fmt.Println("\nShutdown signal received. Closing all connections...")

		server.mutex.Lock()
		for addr, user := range userCache.Users {
			fmt.Printf("Closing connection with %s\n", addr)
			(*user.Conn).Write(protocol.NewMessage(protocol.Die, "Server is shutting down").ToBytes())
			(*user.Conn).Close()
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

		server.cache.Add(conn.RemoteAddr().String(), &state.User{Conn: &conn, Name: ""})
		go server.Handle(conn)
	}
}

func (s *Server) Handle(conn net.Conn) {
	defer func() {
		s.cache.Delete(conn.RemoteAddr().String())
		conn.Close()
	}()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil || n == 0 {
			if s.shutdown {
				return
			}
			s.cache.Delete(conn.RemoteAddr().String())
			return
		}

		m := protocol.ParseMessage(buffer[:n])

		switch m.Type {
		case protocol.SetName:
			fmt.Printf("%s connected as %s\n", conn.RemoteAddr().String(), m.Payload)
			s.cache.Add(conn.RemoteAddr().String(), &state.User{Conn: &conn, Name: m.Payload})

		case protocol.SendMessage:
			if user, ok := s.cache.Get(conn.RemoteAddr().String()); ok {
				fmt.Println(user.Name + ": " + m.Payload)
			}

		case protocol.ListUsers:
			fmt.Println("List of users:")
			for addr, user := range s.cache.Users {
				fmt.Printf("- %s (%s)\n", user.Name, addr)
			}

		case protocol.Logout:
			fmt.Printf("%s disconnected\n", conn.RemoteAddr().String())
			s.cache.Delete(conn.RemoteAddr().String())
			return

		case protocol.KickUser:
			for addr, user := range s.cache.Users {
				if user.Name == m.Payload {
					fmt.Printf("Kicking %s\n", addr)
					(*user.Conn).Write(protocol.NewMessage(protocol.Die, "You were kicked").ToBytes())
					(*user.Conn).Close()
					s.cache.Delete(addr)
				}
			}
		}
	}
}
