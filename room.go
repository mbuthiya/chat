package main

import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"

)

type room struct{
	// Forward is a channel that holds incoming messages
	//that should be forwarded to the other clients

	forward chan []byte

	// Join is a channel for clients wishing to join the channel
	join chan *client

	//Leave is a channel for clients wishing to leave the room
	leave chan *client
	
	// join and leave only exist to safely remove clients from the client map

	// clients hold all the current clients in the room
	clients map[*client]bool
}


func newRoom() *room{

	return &room{
		forward: make(chan []byte),
		join: make(chan *client),
		leave: make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run(){
	for{
		// We use select whenever we need to synchronize shared memory
		// or take different actions depending on various activities within our channels
		select{
		case client := <-r.join:
			//joining
			r.clients[client] = true
		case client := <- r.leave:
			//leaving
			delete(r.clients,client) // removing items from a maps
			close(client.send)
		case msg := <- r.forward:
			//forward message to all clients
			for client := range r.clients{
				client.send <- msg
			}
		}

	}

	//The code will watch the 3 channels in our room and the select will 
	//run the code  for that particular case
}

const(
	socketBufferSize = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize:socketBufferSize,WriteBufferSize:socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter,req *http.Request){

	socket, err :=upgrader.Upgrade(w,req,nil)
	if err!=nil{
		log.Fatal("ServeHTTP",err)
		return
	}

	client := &client{
		socket:socket,
		send: make(chan []byte, messageBufferSize),
		room: r,
	}
	r.join <- client
	
	defer func(){r.leave <- client}()
	go client.write()
	client.read()
}