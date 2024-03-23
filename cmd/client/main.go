package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/NicolasLopes7/tcp-chat/cli"
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

	go withHealthCheck(&conn, signalChan)
	go withDuplexConn(&conn, signalChan)

	name, err := cli.GetNameScreen(scanner)

	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = conn.Write(protocol.NewMessage(protocol.SetName, *name).ToBytes())
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	inputChan := make(chan string)

	go func() {
		fmt.Printf("%s >: ", time.Now().Format(time.Kitchen))
		for scanner.Scan() {
			fmt.Printf("%s >: ", time.Now().Format(time.Kitchen))
			inputChan <- scanner.Text()
		}
	}()

	for {
		select {
		case <-signalChan:
			cancel(&conn, signalChan)
		case line := <-inputChan:
			_, err = conn.Write(GetCommandOrMessage(line).ToBytes())
			if err != nil {
				fmt.Println(err)
				return
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
