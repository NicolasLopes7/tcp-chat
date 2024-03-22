package state

import (
	"net"
	"sync"
)

type Client struct {
	Conn *net.Conn
	Name string
}
type ClientStore struct {
	Clients map[string]*Client
	mutex   sync.RWMutex
}

func NewClientStore() *ClientStore {
	return &ClientStore{
		Clients: make(map[string]*Client),
	}
}

func (cs *ClientStore) Add(key string, value *Client) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	cs.Clients[key] = value
}

func (cs *ClientStore) Range(f func(addr string, client *Client) bool) {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()
	for addr, client := range cs.Clients {
		if !f(addr, client) {
			break
		}
	}
}

func (cs *ClientStore) Delete(key string) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	delete(cs.Clients, key)
}

func (cs *ClientStore) Get(key string) (*Client, bool) {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()
	val, ok := cs.Clients[key]
	return val, ok
}
