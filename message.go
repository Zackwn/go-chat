package main

// message types
const (
	SEND_MESSAGE messageType = iota
	JOIN_ROOM
	QUIT_ROOM
	ROOM_MEMBERS
)

type messageType int

type message struct {
	Type messageType `json:"type"`
	Data interface{} `json:"data,omitempty"`
}
