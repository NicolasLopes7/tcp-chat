package services

import (
	"net"

	"github.com/NicolasLopes7/tcp-chat/state"
)

type AbstractUserService interface {
	SetName(name string)
	SetConn(conn *net.Conn)
	GetName() string
}

type UserService struct {
	User *state.User
}

func (us *UserService) SetName(name string) {
	us.User.Name = name
}

func (us *UserService) SetConn(conn *net.Conn) {
	us.User.Conn = conn
}

func (us *UserService) GetName() string {
	return us.User.Name
}
