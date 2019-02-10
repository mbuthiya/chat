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


// read method allows us to read from the client socket via the ReadMessage method

func (c *client) read(){
	defer c.socket.Close()
	for{
		_,msg,err := c.socket.ReadMessage()

		// if the Read message encounters an error we break the loop
		if err != nil{
			return
		}
		//Messages from the client are sent to the rooms forward channel
		c.room.forward <- msg
	}
}

// The defer keyword to tidy up and close the web socket connection when the function returns

// write method allows us to write to the client socket with data from the send channel

func (c *client) write(){
	defer c.socket.Close()
	
	for msg := range c.send{
		err := c.socket.WriteMessage(websocket.TextMessage,msg)

		if err!= nil{
			return
		}
	}
}