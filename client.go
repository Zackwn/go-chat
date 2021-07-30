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
	client.ID = fmt.Sprint(rand.Int())

	log.Printf("new client: %v\n", client.ID)
	return client
}

type client struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	conn      *websocket.Conn
	room      *room
	join_room chan<- joinRoomInfo
	send      chan message
}

func (c *client) readPump() {
	for {
		var msg message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			c.quitCurrentRoom()
			c.conn.Close()
			log.Println(err)
			return
		}
		switch msg.Type {
		case SEND_MESSAGE:
			text, ok := msg.Data.(string)
			if !ok {
				log.Println("invalid data from client: ", msg.Data)
			}
			c.room.broadcast <- message{
				Type: SEND_MESSAGE,
				Data: struct {
					Client  *client `json:"client"`
					Content string  `json:"content"`
				}{
					Client:  c,
					Content: text,
				},
			}

		case QUIT_ROOM:
			c.quitCurrentRoom()

		case JOIN_ROOM:
			roomName, ok := msg.Data.(string)
			if !ok {
				log.Println("invalid message from client: ", msg.Data)
				return
			}
			c.quitCurrentRoom()
			c.join_room <- joinRoomInfo{client: c, roomName: roomName}
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
		client.room.unregister <- client
		client.room.broadcast <- message{
			Type: QUIT_ROOM,
			Data: client,
		}
		client.room = nil
	}
}
