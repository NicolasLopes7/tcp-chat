package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NicolasLopes7/tcp-chat/protocol"
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

	go func() {
		t := time.NewTicker(1 * time.Second)
		defer t.Stop()

		for range t.C {
			err := checkServerStatus(&conn)
			if err != nil {
				fmt.Print("\n\nServer is down\n")
				cancel(&conn, signalChan)
			}
		}
	}()

	go func() {
		for {
			message, err := protocol.ReadMessage(&conn)
			if err != nil {
				fmt.Println("Error reading message:", err)
				return
			}
			if message.Type == protocol.Die {
				fmt.Println(message.Payload)
				cancel(&conn, signalChan)
				return
			}
		}
	}()

	fmt.Println("Enter your name: ")
	scanner.Scan()
	name := scanner.Text()

	_, err = conn.Write(protocol.NewMessage(protocol.SetName, name).ToBytes())
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	inputChan := make(chan string)

	go func() {
		fmt.Printf(">: ")
		for scanner.Scan() {
			fmt.Printf(">: ")
			inputChan <- scanner.Text()
		}
	}()

	for {
		select {
		case <-signalChan:
			cancel(&conn, signalChan)
		case line := <-inputChan:
			_, err = conn.Write(protocol.NewMessage(protocol.SendMessage, line).ToBytes())
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func cancel(conn *net.Conn, signalChan chan os.Signal) {
	(*conn).Write(protocol.NewMessage(protocol.Logout, "disconnect").ToBytes())
	close(signalChan)
	(*conn).Close()
	os.Exit(0)
}

func checkServerStatus(conn *net.Conn) error {
	_, err := (*conn).Write(protocol.NewMessage(protocol.Ping, "").ToBytes())
	if err != nil {
		return err
	}

	return nil
}
