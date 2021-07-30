package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

func NewClient(conn *websocket.Conn) *client {
	rand.Seed(time.Now().Unix())

	client := new(client)
	client.send = make(chan message)
	client.conn = conn
	client.id = fmt.Sprint(rand.Int())

	log.Printf("new client: %v\n", client.id)
	return client
}

type client struct {
	id        string
	name      string
	conn      *websocket.Conn
	room      *room
	join_room chan<- joinRoomInfo
	send      chan message
}

func (client *client) readPump() {
	for {
		var msg message
		err := client.conn.ReadJSON(&msg)
		if err != nil {
			client.quitCurrentRoom()
			client.conn.Close()
			log.Println(err)
			return
		}
		switch msg.Type {
		case SEND_MESSAGE:
			text, ok := msg.Data.(string)
			if !ok {
				log.Println("invalid data from client: ", msg.Data)
			}
			client.room.broadcast <- message{
				Type: SEND_MESSAGE,
				Data: struct {
					Name    string `json:"name"`
					Content string `json:"content"`
				}{
					Name:    client.name,
					Content: text,
				},
			}

		case QUIT_ROOM:
			client.quitCurrentRoom()

		case JOIN_ROOM:
			roomName, ok := msg.Data.(string)
			if !ok {
				log.Println("invalid message from client: ", msg.Data)
				return
			}
			client.join_room <- joinRoomInfo{client: client, roomName: roomName}
		}
	}
}

func (client *client) writePump() {
	for {
		message, open := <-client.send
		if !open {
			return
		}
		json, err := json.Marshal(&message)
		if err != nil {
			log.Println(err)
			return
		}
		client.conn.WriteMessage(websocket.TextMessage, json)
	}
}

func (client *client) quitCurrentRoom() {
	if client.room != nil {
		delete(client.room.members, client.id)
		client.room.broadcast <- message{
			Type: QUIT_ROOM,
			Data: client.name,
		}
		client.room = nil
	}
}
