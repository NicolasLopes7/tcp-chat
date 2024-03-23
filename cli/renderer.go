package cli

import (
	"bufio"
	"fmt"

	"github.com/NicolasLopes7/tcp-chat/services"
)

type RendererContainer struct {
	roomService   services.AbstractRoomService
	userService   services.AbstractUserService
	stdinConsumer services.AbstractStdinConsumer
	scanner       *bufio.Scanner
}

func NewPageRender(container *RendererContainer) func(name string) error {
	return func(name string) error {
		switch name {
		case "get-name":
			name, err := GetNameScreen(container.scanner)
			if err != nil {
				fmt.Println("Error: ", err)
				return err
			}

			container.userService.SetName(*name)
			return nil
		case "select-room":
			room, err := SelectRoomScreen(container.scanner, container.roomService.GetRooms(container.userService.GetName()))
			if err != nil {
				fmt.Println("Error: ", err)
				return err
			}

			err = container.roomService.SubscribeToRoom(*room)
			if err != nil {
				return err
			}
			return nil
		case "chat":
			inputChan := ChatScreen(container.scanner)
			err := container.stdinConsumer.Parse(inputChan)

			if err != nil {
				fmt.Println("Error: ", err)
				return err
			}

			return nil
		}
		return nil
	}
}

func NewRenderer(container *RendererContainer) {
	pages := []*Page{
		{name: "get-name"},
		{name: "select-room"},
		{name: "chat"},
	}

	history := NewHistory(pages)
	render := NewPageRender(container)

	for {
		currentPage := history.pages[history.currentIdx]

		err := render(currentPage.name)
		if err != nil {
			render(currentPage.name)
			continue
		}

		err = history.Forward()
		if err != nil {
			break
		}
	}
}
