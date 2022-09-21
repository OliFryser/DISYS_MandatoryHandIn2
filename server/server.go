package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/OliFryser/DISYS_MandatoryHandIn2/TCP"

	"google.golang.org/grpc"
)

type Server struct {
	TCP.UnimplementedTransmitDataServer
}

func (s *Server) GetTransmitData(ctx context.Context, in *TCP.ClientHandshake) (*TCP.ServerHandshake, error) {
	fmt.Printf("Received XXX request")
	return &TCP.ServerHandshake{Reply: "Your reply here"}, nil
}

func main() {
	// Create listener tcp on port 9080
	list, err := net.Listen("tcp", ":9080")
	if err != nil {
		log.Fatalf("Failed to listen on port 9080: %v", err)
	}
	grpcServer := grpc.NewServer()
	TCP.RegisterTransmitDataServer(grpcServer, &Server{})

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to server %v", err)
	}
}
