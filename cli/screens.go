package cli

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/NicolasLopes7/tcp-chat/state"
)

func GetNameScreen(scanner *bufio.Scanner) (*string, error) {
	fmt.Println("Enter your name: ")
	scanner.Scan()
	name := scanner.Text()

	if strings.Contains(name, "*") || name == "" {
		return nil, fmt.Errorf("invalid name. Names are required, and must not contain *")
	}
	return &name, nil
}

func SelectRoomScreen(scanner *bufio.Scanner, rooms []*state.Room) (*string, error) {
	fmt.Println("Select the room you want to join: ")

	for idx, room := range rooms {
		fmt.Printf("%d) #%s\n", idx, room.Name)
	}

	fmt.Println(">: ")
	scanner.Scan()
	room := scanner.Text()

	idx, err := strconv.Atoi(room)

	if err != nil || idx > len(rooms) || idx < 1 {
		return nil, fmt.Errorf("invalid room number")
	}

	return &rooms[idx].Name, nil
}

func ChatScreen(scanner *bufio.Scanner) chan string {
	inputChan := make(chan string)

	fmt.Printf("%s >: ", time.Now().Format(time.Kitchen))
	for scanner.Scan() {
		fmt.Printf("%s >: ", time.Now().Format(time.Kitchen))
		inputChan <- scanner.Text()
	}

	return inputChan
}
