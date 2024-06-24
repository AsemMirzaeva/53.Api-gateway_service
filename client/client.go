package main

import (
	"bufio"
	chatpb "chat/proto"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	conn, err := grpc.NewClient("localhost:8001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := chatpb.NewChatServiceClient(conn)
	stream, err := client.Chat(context.Background())
	if err != nil {
		log.Fatalf("Failed to create a stream: %v", err)
	}

	go func() {
		for {
			msg, err := stream.Recv()
			st := status.Convert(err)
			if st.Code() == codes.Unavailable {
				if err := stream.CloseSend(); err != nil {
					log.Fatalf("Couldn't connect to the server: %v", st.Code())
				}
			}

			if err != nil {
				log.Fatalf("Failed to receive message: %v", err)
			}
			fmt.Printf("%s [%s] (%s): %s\n", time.Unix(msg.Timestamp, 0).Format(time.DateTime), msg.User, msg.IpAddress, msg.Message)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		err := stream.Send(&chatpb.ChatMessage{
			User:      "TestUser",
			Message:   message,
			Timestamp: time.Now().Unix(),
		})

		if err != nil {
			log.Fatalf("Failed to send a message: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading from input: %v", err)
	}
}
