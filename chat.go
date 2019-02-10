package main

import(
	"github.com/gorilla/websocket"
)

// Client represents a single chatting user

type client struct{
	
	socket *websocket.Conn // Socket is the web socket for the client

	// Send is a channel on which messages are sent
	send chan []byte

	// Room is the room this client is chatting in

	room *room

}