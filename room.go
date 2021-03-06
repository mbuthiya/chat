package main

import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	"github.com/mbuthiya/tracer"

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

	// Tracer will recieve trace information activity in the room
	tracer tracer.Tracer
}


func newRoom() *room{

	return &room{
		forward: make(chan []byte),
		join: make(chan *client),
		leave: make(chan *client),
		clients: make(map[*client]bool),
		tracer: tracer.Off(),
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
			r.tracer.Trace("New Client has Joined")
		case client := <- r.leave:
			//leaving
			delete(r.clients,client) // removing items from a maps
			close(client.send)

			r.tracer.Trace("Client left")
		case msg := <- r.forward:

			r.tracer.Trace("Message recieced: ",string(msg))
			//forward message to all clients
			for client := range r.clients{
				client.send <- msg
				r.tracer.Trace("--sent to client")
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


// To use websockets we need to upgrade tge HTTP connection using the websocket.Upgrader type
var upgrader = &websocket.Upgrader{ReadBufferSize:socketBufferSize,WriteBufferSize:socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter,req *http.Request){

	// We get a web socket by calling the upgrader.Upgrade method
	socket, err :=upgrader.Upgrade(w,req,nil)
	if err!=nil{
		log.Fatal("ServeHTTP",err)
		return
	}

	// Create the client object 
	client := &client{
		socket:socket,
		send: make(chan []byte, messageBufferSize),
		room: r,
	}

	// Pass the client to the room 
	r.join <- client
	
	// When the client is finished we can remove the client from the room
	defer func(){r.leave <- client}()
	
	//we create a go routine to write to the web socket in the background
	go client.write()

	// Call the read method to keep the connection open
	client.read()
}