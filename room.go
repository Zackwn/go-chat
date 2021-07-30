package main

import (
	"encoding/json"
)

func NewRoom(members map[string]*client) *room {
	room := new(room)
	room.members = members
	room.broadcast = make(chan message)
	room.register = make(chan *client)
	room.unregister = make(chan *client)
	return room
}

type room struct {
	// map[id]client
	members    members
	broadcast  chan message
	register   chan *client
	unregister chan *client
}

type members map[string]*client

func (room *room) run() {
	for {
		select {
		case msg := <-room.broadcast:
			for memberid := range room.members {
				room.members[memberid].send <- msg
			}
		case client := <-room.register:
			room.members[client.ID] = client
		case client := <-room.unregister:
			delete(room.members, client.ID)
		}
	}
}

func (members members) MarshalJSON() ([]byte, error) {
	out := make([]client, len(members))
	i := 0
	for id := range members {
		out[i] = *members[id]
		i++
	}
	return json.Marshal(out)
}
