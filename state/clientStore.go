package state

import "sync"

type ClientStore struct {
	clients map[string]string
	mutex   sync.RWMutex
}

func NewClientStore() *ClientStore {
	return &ClientStore{
		clients: make(map[string]string),
	}
}

func (cs *ClientStore) Add(key string, value string) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	cs.clients[key] = value
}

func (cs *ClientStore) Delete(key string) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	delete(cs.clients, key)
}

func (cs *ClientStore) Get(key string) (string, bool) {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()
	val, ok := cs.clients[key]
	return val, ok
}
