package main

import (
	"log"
	"msqrd/pkg/api"
	"msqrd/pkg/chat"
	"net"

	"google.golang.org/grpc"
)

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	api.RegisterChatServiceServer(s, new(chat.GRPCServer))

	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}
