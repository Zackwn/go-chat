package main

import (
	"log"

	"github.com/gorilla/websocket"
)

func NewServer() *server {
	joinRoom := make(chan joinRoomInfo)
	rooms := make(map[string]*room)
	return &server{rooms: rooms, join_room: joinRoom}
}

type server struct {
	// map[room_name]room
	rooms     map[string]*room
	join_room chan joinRoomInfo
}

type joinRoomInfo struct {
	client   *client
	roomName string
}

func (server *server) run() {
	for {
		info := <-server.join_room
		room, exists := server.rooms[info.roomName]
		// create room if it don't exists
		if !exists {
			room = NewRoom(make(map[string]*client))
			go room.run()
			// register room
			server.rooms[info.roomName] = room
		}
		// register client
		room.register <- info.client
		// inform room members
		info.client.room = room
		info.client.room.broadcast <- message{
			Type: JOIN_ROOM,
			Data: info.client.name,
		}
	}
}

func (server *server) newClient(conn *websocket.Conn) {
	var r map[string]string
	err := conn.ReadJSON(&r)
	if err != nil {
		log.Println(err)
		return
	}
	client := NewClient(conn)
	client.join_room = server.join_room
	client.name = r["name"]

	go client.readPump()
	go client.writePump()
}
