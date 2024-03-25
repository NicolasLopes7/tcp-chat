package cli

import (
	"bufio"
	"fmt"

	"github.com/NicolasLopes7/tcp-chat/services"
)

type RendererContainer struct {
	RoomService   services.AbstractRoomService
	UserService   services.AbstractUserService
	StdinConsumer services.AbstractStdinConsumer
	Scanner       *bufio.Scanner
}

func NewPageRender(container *RendererContainer) func(name string) error {
	return func(name string) error {
		switch name {
		case "get-name":
			name, err := GetNameScreen(container.Scanner)
			if err != nil {
				fmt.Println("Error: ", err)
				return err
			}

			container.UserService.SetName(*name)
			return nil
		case "select-room":
			room, err := SelectRoomScreen(container.Scanner, container.RoomService.GetRooms(container.UserService.GetName()))
			if err != nil {
				fmt.Println("Error: ", err)
				return err
			}

			err = container.RoomService.SubscribeToRoom(*room)
			if err != nil {
				return err
			}
			return nil
		case "chat":
			inputChan := ChatScreen(container.Scanner)
			err := container.StdinConsumer.Parse(inputChan)

			if err != nil {
				fmt.Println("Error: ", err)
				return err
			}

			return nil
		}
		return nil
	}
}

func Render(container *RendererContainer) {
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
