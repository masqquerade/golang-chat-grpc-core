package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"msqrd/pkg/api"
	"sync"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	c := api.NewChatServiceClient(conn)

	stream, err := c.Connect(context.Background(), &api.Null{})
	if err != nil {
		log.Fatalf("stream connection error: %s", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		for {
			defer wg.Done()

			msg, err := stream.Recv()
			if err != nil && err != io.EOF {
				fmt.Printf("reading message error: %s", err)
				return
			}

			if msg != nil {
				fmt.Println(msg)
			}
		}
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			msg := &api.Message{
				Uname: "user",
				Msg:   "hi!!!",
			}

			_, err := c.SendMessage(context.Background(), msg)

			if err != nil {
				fmt.Printf("sending message error: %s", err)
				return
			}
		}
	}()

	wg.Wait()
}
