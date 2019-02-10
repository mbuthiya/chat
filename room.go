package main

type room struct{
	// Forward is a channel that holds incoming messages
	//that should be forwarded to the other clients

	forward chan []byte
}