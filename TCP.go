package main

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"time"
)

type Packet struct {
	syn  int
	ack  int
	hash uint32
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
	time.Sleep(10000 * time.Millisecond)
}

func serverThread(s *Server) {
	//The server receives the first handshake
	message := <-s.reciever
	fmt.Printf("Server received SYN: %d \n", message.syn)
	//The server takes the syn number and adds 1, then generates a new random integer
	message.ack = message.syn + 1
	message.syn = rand.Intn(420)
	syn := message.syn

	//The server sends the second handshake
	s.sender <- message

	//The server receives the third handshake
	message = <-s.reciever

	//Check if the acknowledgement is correct
	if message.ack != syn+1 {
		fmt.Println("!!!Wrong SYN number!!!")
	}
	fmt.Printf("Server received ACK: %d \n", message.ack)

	var lastsyn int
	outOfSync := false
	var lastCorrectSyn int

	for {
		select {
		//Server receives data
		case message := <-s.reciever:
			fmt.Printf("Server received data: %s \n", message.data)

			lastsyn = syn
			syn = message.syn

			//Check message reordering and data loss
			if (syn != 0 && lastsyn != syn-1) || (message.hash != hash(message.data)) {
				if !outOfSync {
					lastCorrectSyn = lastsyn
				}
				outOfSync = true
			}

			//Sleeping to avoid race conditions
			time.Sleep(100 * time.Millisecond)
		//empties the sender channel
		case <-s.sender:

		//When the client is finished, return a packet with ACK
		default:
			message := new(Packet)
			if outOfSync {
				message.ack = lastCorrectSyn
				outOfSync = !outOfSync
			} else {
				message.ack = syn + 1
			}
			s.sender <- *message
		}

	}

}

func clientThread(c *Client) {
	var message Packet
	message.syn = rand.Intn(420)
	syn := message.syn

	//Send first handshake
	c.sender <- message

	//Receive second handshake
	message = <-c.reciever

	//Check server handshake
	if message.ack != syn+1 {
		fmt.Println("Wrong SYN number")
	}
	fmt.Printf("Client received ACK: %d\n", message.ack)
	fmt.Printf("Client received SYN: %d\n", message.syn)

	//Send third handshake
	message.ack = message.syn + 1
	c.sender <- message

	data := ""

	transferSuccess := false

	//Transfer data
	for !transferSuccess {
		for i := 0; i < 10; i++ {

			message := new(Packet)
			data += "HELLO "
			message.data = data
			message.hash = hash(data)
			//If one wants to see the error handling, uncomment the below code and an error will occur
			/*if i == 5 {
				message.syn = 1000
			} else {
				message.syn = i
			}*/
			message.syn = i

			syn = message.syn

			c.sender <- *message
		}

		message = <-c.reciever
		fmt.Printf("Client received ACK after data transfer: %d \n", message.ack)

		if message.ack == syn+1 {
			transferSuccess = true
		} else {
			fmt.Printf("DATA OUT OF ORDER OR CORRUPTED!\nPlease try again.")
		}
	}
	fmt.Printf("Data transfer succesfull!\n")
	fmt.Printf("Program will exit soon\n")
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
