package main

type room struct{
	// Forward is a channel that holds incoming messages
	//that should be forwarded to the other clients

	forward chan []byte

	// Join is a channel for clients wishing to join the channel
	join chan *client

	//Leave is a channel for clients wishing to leave the room
	leave chan *client
	
	// clients hold all the current clients in the room

	client map[*client]bool
}