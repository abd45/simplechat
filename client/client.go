package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	pb "github.com/abd45/simplechat/simplechat"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	"time"
)

var (
	serverAddr = flag.String("server-address", "localhost:10001", "The server address in the format of host:port")
	sender     = flag.String("sender", "NOUSER137#2310945", "The username of the sender who is sending the chat")
	receiver   = flag.String("receiver", "NOUSER137#2310945", "The username of the person for whom this chat is intended to.")
)

func sendMessage(ctx context.Context, client pb.SimpleChatClient, username string) {
	reader := bufio.NewReader(os.Stdin)
	var stream pb.SimpleChat_SendMessageClient
	for {
		text, _ := reader.ReadString('\n')
		stream, err := client.SendMessage(ctx)
		if err != nil {
			log.Fatalf("%v.SendMessage(_) = _, %v", client, err)
		}

		if err := stream.Send(&pb.Conversation{Username: username, Ping: text}); err != nil {
			log.Fatalf("%v.Send(%v) = %v", stream, "", err)
		}
	}
	_, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}
}

func receiveMessage(ctx context.Context, client pb.SimpleChatClient, user *pb.User) {
	stream, err := client.ReceiveMessage(ctx, user)
	if err != nil {
		log.Fatalf("%v.ReceiveMessage(_) = _, %v", client, err)
	}

	for {
		ping, err := stream.Recv()
		if err == io.EOF {
			return
		}

		if err != nil {
			log.Fatalf("Failed to receive a ping : %v", err)
		}

		fmt.Println("-> ", ping.GetPing())
	}
}
func main() {
	flag.Parse()

	if *sender == "NOUSER137#2310945" {
		log.Fatal("Please set a username for yourself with --sender flag")
	}

	if *receiver == "NOUSER137#2310945" {
		log.Fatal("Please set a username for receiver with --receiver flag")
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	defer conn.Close()
	client := pb.NewSimpleChatClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	_, err = client.RegisterUser(ctx, &pb.User{Username: *sender})
	if err != nil {
		log.Fatal(err)
	}

	waitc := make(chan struct{})
	go func() {
		for {
			receiveMessage(ctx, client, &pb.User{Username: *sender})
			time.Sleep(100 * time.Millisecond)
		}
	}()
	go sendMessage(ctx, client, *receiver)
	<-waitc
}
