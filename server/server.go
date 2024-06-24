package main

import (
	"chat/pq"
	chatpb "chat/proto"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type server struct {
	chatpb.UnimplementedChatServiceServer
	mu      sync.Mutex
	clients map[string]chatpb.ChatService_ChatServer
	chatDB  *pq.ChatDB
}

func newServer(chatDB *pq.ChatDB) *server {
	return &server{
		clients: make(map[string]chatpb.ChatService_ChatServer),
		chatDB:  chatDB,
	}
}

func (s *server) Chat(stream chatpb.ChatService_ChatServer) error {
	p, ok := peer.FromContext(stream.Context())
	if !ok {
		return status.Error(codes.Internal, "could not get peer info")
	}

	clientID := p.Addr.String()
	network := p.LocalAddr.Network()
	fmt.Printf("connected [%s]:%s\n", network, clientID)

	s.mu.Lock()
	s.clients[clientID] = stream
	s.mu.Unlock()

	messages, err := s.chatDB.LoadMessages()
	if err != nil {
		return status.Error(codes.Internal, "failed to load messages from database")
	}
	for _, msg := range messages {
		if err := stream.Send(msg); err != nil {
			return err
		}
	}

	defer func() {
		s.mu.Lock()
		delete(s.clients, clientID)
		s.mu.Unlock()
	}()

	for {
		msg, err := stream.Recv()
		if err != nil {
			return err
		}

		msg.IpAddress = clientID

		if err := s.chatDB.SaveMessage(msg); err != nil {
			return status.Error(codes.Internal, "failed to save message to database")
		}

		s.broadcastMessage(msg, clientID)
	}
}

func (s *server) broadcastMessage(msg *chatpb.ChatMessage, senderID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, client := range s.clients {
		if id == senderID {
			continue
		}

		if err := client.Send(msg); err != nil {
			log.Printf("Error sending message to client %v: %v", id, err)
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":8001")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	db, err := pq.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to  database: %v", err)
	}
	chatDB := pq.NewChatDB(db)

	s := grpc.NewServer()
	chatpb.RegisterChatServiceServer(s, newServer(chatDB))

	log.Printf("Server is listening on port 8001")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
