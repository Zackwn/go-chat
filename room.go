package main

func NewRoom(members map[string]*client) *room {
	room := new(room)
	room.members = members
	room.broadcast = make(chan message)
	room.register = make(chan *client)
	return room
}

type room struct {
	// map[id]client
	members   map[string]*client
	broadcast chan message
	register  chan *client
}

func (room *room) run() {
	for {
		select {
		case msg := <-room.broadcast:
			for memberid := range room.members {
				room.members[memberid].send <- msg
			}
		case client := <-room.register:
			room.members[client.id] = client
		}
	}
}
