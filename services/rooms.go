package services

import (
	"net"

	"github.com/NicolasLopes7/tcp-chat/protocol"
	"github.com/NicolasLopes7/tcp-chat/state"
)

type AbstractRoomService interface {
	GetRooms(userName string) []*state.Room
	SubscribeToRoom(room string) error
	UnsubscribeToRoom(room string) error
}

type RoomService struct {
	rooms []*state.Room
	conn  *net.Conn
}

func (rs *RoomService) GetRooms(userName string) []*state.Room {
	var result []*state.Room
	for _, room := range rs.rooms {
		if room.Acl == nil {
			continue
		}

		if room.Acl.Public {
			result = append(result, room)
		}

		for _, user := range room.Acl.Users {
			if user.Name == userName {
				result = append(result, room)
				break
			}
		}
	}

	return result
}

func (rs *RoomService) SubscribeToRoom(room string) error {
	_, err := (*rs.conn).Write(protocol.NewMessage(protocol.SubscribeToRoom, room).ToBytes())

	if err != nil {
		return err
	}

	return nil
}

func (rs *RoomService) UnsubscribeToRoom(room string) error {
	_, err := (*rs.conn).Write(protocol.NewMessage(protocol.UnsubscribeToRoom, room).ToBytes())

	if err != nil {
		return err
	}

	return nil
}
