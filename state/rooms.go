package state

type Acl struct {
	Public bool
	Users  []*User
}

type Room struct {
	Name string
	Acl  *Acl
}
