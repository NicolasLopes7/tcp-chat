package state

import (
	"sync"
)

type UserCache struct {
	Users map[string]*User
	mutex sync.RWMutex
}

func NewUserCache() *UserCache {
	return &UserCache{
		Users: make(map[string]*User),
	}
}

func (cs *UserCache) Add(key string, value *User) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	cs.Users[key] = value
}

func (cs *UserCache) Delete(key string) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	delete(cs.Users, key)
}

func (cs *UserCache) Get(key string) (*User, bool) {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()
	val, ok := cs.Users[key]
	return val, ok
}
