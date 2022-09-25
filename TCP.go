package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Packet struct {
	syn  int
	ack  int
	data string
}

type Server struct {
	sender   chan Packet
	reciever chan Packet
}

type Client struct {
	sender   chan Packet
	reciever chan Packet
}

func main() {
	var server Server
	var client Client
	server.sender = make(chan Packet)
	client.sender = make(chan Packet)
	server.reciever = client.sender
	client.reciever = server.sender

	go serverThread(&server)
	go clientThread(&client)
	time.Sleep(1000 * time.Millisecond)
}

func serverThread(s *Server) {
	message := <-s.reciever
	fmt.Printf("Server received SYN: %d \n", message.syn)
	message.ack = message.syn + 1
	message.syn = rand.Intn(420)
	syn := message.syn

	s.sender <- message
	message = <-s.reciever

	if message.ack != syn+1 {
		fmt.Println("!!!Wrong SYN number!!!")
	}
	fmt.Printf("Server received ACK: %d \n", message.ack)

	for {
		select {
		case message = <-s.reciever:
			fmt.Printf("Server received data: %s \n", message.data)
		default:
			message.ack = message.syn + 1
			s.sender <- message
		}

	}

}

func clientThread(c *Client) {
	var message Packet
	message.syn = rand.Intn(420)
	syn := message.syn

	c.sender <- message
	message = <-c.reciever

	if message.ack != syn+1 {
		fmt.Println("Wrong SYN number")
	}
	fmt.Printf("Client received ACK: %d\n", message.ack)
	fmt.Printf("Client received SYN: %d\n", message.syn)

	message.ack = message.syn + 1
	c.sender <- message

	data := ""

	//data
	for i := 0; i < 10; i++ {
		data += "HELLO "
		message.data = data
		message.syn = i
		syn = message.syn
		c.sender <- message
	}

	message = <-c.reciever

	fmt.Printf("Client received ACK after data transfer: %d \n", message.ack)
	if message.ack != syn+1 {
		fmt.Println("Wrong SYN number")
	}

	fmt.Printf("Data transfer succesful!")
}
