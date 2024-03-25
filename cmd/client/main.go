package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NicolasLopes7/tcp-chat/cli"
	"github.com/NicolasLopes7/tcp-chat/network"
	"github.com/NicolasLopes7/tcp-chat/protocol"
	"github.com/NicolasLopes7/tcp-chat/services"
	"github.com/NicolasLopes7/tcp-chat/state"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	scanner := bufio.NewScanner(os.Stdin)

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	roomService := &services.RoomService{
		Rooms: []*state.Room{},
		Conn:  &conn,
	}

	userService := &services.UserService{
		User: &state.User{Conn: &conn},
	}

	writerService := &network.WriterService{
		Conn: &conn,
	}

	stdinConsumer := &services.StdinConsumer{
		WriterService: writerService,
	}

	cli.Render(&cli.RendererContainer{
		RoomService:   roomService,
		UserService:   userService,
		Scanner:       scanner,
		StdinConsumer: stdinConsumer,
	})

	go withHealthCheck(&conn, signalChan)
	go withDuplexConn(&conn, signalChan)
	go WithSignals(&conn, signalChan)

}

func cancel(conn *net.Conn, signalChan chan os.Signal) {
	(*conn).Write(protocol.NewMessage(protocol.Logout, "disconnect").ToBytes())
	close(signalChan)
	(*conn).Close()
	os.Exit(0)
}

func withDuplexConn(conn *net.Conn, c chan os.Signal) {
	for {
		message, err := protocol.ReadMessage(conn)
		if err != nil {
			fmt.Println("Error reading message:", err)
			return
		}
		fmt.Println(message.Payload)
		if message.Type == protocol.Die {
			cancel(conn, c)
			return
		}
	}
}
func withHealthCheck(conn *net.Conn, c chan os.Signal) {
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()

	for range t.C {
		err := checkServerStatus(conn)
		if err != nil {
			fmt.Print("\n\nServer is down\n")
			cancel(conn, c)
		}
	}
}
func checkServerStatus(conn *net.Conn) error {
	_, err := (*conn).Write(protocol.NewMessage(protocol.Ping, "").ToBytes())
	if err != nil {
		return err
	}

	return nil
}

func WithSignals(conn *net.Conn, c chan os.Signal) {
	for {
		select {
		case <-c:
			cancel(conn, c)
		}
	}
}
