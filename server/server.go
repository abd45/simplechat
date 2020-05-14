package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/abd45/simplechat/simplechat"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"sync"
)

var (
	serverAddr = flag.String("address", "localhost:10001", "The server address in the format of host:port")
	conver     *pb.Conversation
)

type Data struct {
	username string
	message  []string
}

type simpleChatServer struct {
	pb.UnimplementedSimpleChatServer
	users []*Data
	mu    sync.Mutex
}

func (s *simpleChatServer) SendMessage(stream pb.SimpleChat_SendMessageServer) error {
	for {
		conv, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.Acknowledge{
				Ack: true,
			})
		}
		if err != nil {
			return err
		}

		noSuchUserExists := false
		s.mu.Lock()
		for _, i := range s.users {
			if i.username == conv.GetUsername() {
				noSuchUserExists = true
				i.message = append(i.message, conv.GetPing())
			}
		}
		s.mu.Unlock()
		if noSuchUserExists {
			return fmt.Errorf("No user named %s found, please ask %s to join the room", conv.GetUsername(), conv.GetUsername())
		}
	}

}

func (s *simpleChatServer) ReceiveMessage(user *pb.User, stream pb.SimpleChat_ReceiveMessageServer) error {
	for _, i := range s.users {
		if i.username == user.GetUsername() {
			if i.message != nil {
				for _, j := range i.message {
					stream.Send(&pb.Conversation{Ping: j, Username: i.username})
				}
				s.mu.Lock()
				i.message = nil
				s.mu.Unlock()
			}
		}
	}
	return nil
}

func (s *simpleChatServer) RegisterUser(ctx context.Context, user *pb.User) (ack *pb.Acknowledge, err error) {
	data := &Data{username: user.GetUsername(), message: make([]string, 0, 10)}
	s.mu.Lock()
	userExists := false
	for _, i := range s.users {
		if i.username == user.GetUsername() {
			userExists = true
		}
	}
	if userExists {
		err = fmt.Errorf("User with username %s exists. Please run the client with different username.", user.GetUsername())
	} else {
		s.users = append(s.users, data)
	}
	s.mu.Unlock()

	if err != nil {
		return nil, err
	}

	fmt.Println("Total users: ", len(s.users))

	ack = &pb.Acknowledge{Ack: true}
	return ack, nil
}

func NewSimpleChatServer() *simpleChatServer {
	return &simpleChatServer{
		users: make([]*Data, 0, 10),
	}
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", *serverAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleChatServer(grpcServer, NewSimpleChatServer())

	fmt.Println("Starting chat server")
	grpcServer.Serve(lis)
}
