package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

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
			conn.Write(protocol.NewMessage(protocol.Logout, "disconnect").ToBytes())
			close(signalChan)
			conn.Close()
			os.Exit(0)
		case line := <-inputChan:
			_, err = conn.Write(protocol.NewMessage(protocol.SendMessage, line).ToBytes())
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
