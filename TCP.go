package main

import (
	"fmt"
	"time"
)

type Server struct {
	sender   chan string
	reciever chan string
}

type Client struct {
	sender   chan string
	reciever chan string
}

func main() {
	var server Server
	var client Client
	server.sender = make(chan string)
	client.sender = make(chan string)
	server.reciever = client.sender
	client.reciever = server.sender

	go serverThread(&server)
	go clientThread(&client)
	time.Sleep(1000 * time.Millisecond)
}

func serverThread(s *Server) {
	message := <-s.reciever
	fmt.Printf("Server recieved message: %s \n", message)
	s.sender <- "Hello Client!"

}

func clientThread(c *Client) {
	c.sender <- "Hello Server"
	message := <-c.reciever
	fmt.Printf("Client recieved message: %s\n", message)
	c.sender <- "Hello again, here's some data\n"

}
