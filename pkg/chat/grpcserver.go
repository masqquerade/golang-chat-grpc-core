package chat

import (
	"context"
	"fmt"
	"msqrd/pkg/api"
	"sync"
)

type GRPCServer struct {
	Streams []api.ChatService_ConnectServer
	Err     chan error
}

func (s *GRPCServer) Connect(_ *api.Null, stream api.ChatService_ConnectServer) error {
	s.Streams = append(s.Streams, stream)
	s.Err = make(chan error)

	return nil
}

func (s *GRPCServer) SendMessage(ctx context.Context, msg *api.Message) (*api.Null, error) {
	wg := sync.WaitGroup{}

	for _, stream := range s.Streams {
		wg.Add(1)

		go func(st api.ChatService_ConnectServer) {
			defer wg.Done()

			if err := st.Send(msg); err != nil {
				s.Err <- err
				return
			}

		}(stream)
	}

	wg.Wait()
	fmt.Println("send")
	return &api.Null{}, <-s.Err
}
