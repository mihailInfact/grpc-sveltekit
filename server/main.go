package main

import (
	"context"
	"greeter/internal/store"
	pb "greeter/pkg/greeter/proto"
	"log"
	"net"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	pb.UnimplementedToDoServiceServer
	Store *store.Store
}

func (s *server) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	todoItem := req.GetItem()

	log.Printf("creating item: %+v\n", todoItem)

	var id string
	var createdAt time.Time

	err := s.Store.QueryRowContext(ctx,
		"INSERT INTO todos (title, description, status) VALUES (?, ?, ?) RETURNING id, created_at",
		todoItem.GetTitle(),
		todoItem.GetDescription(),
		todoItem.GetStatus(),
	).Scan(&id, &createdAt)

	if err != nil {
		log.Printf("Failed to insert todo item: %v", err)
		return nil, err
	}

	// Implement the logic to create a new todo item
	return &pb.CreateResponse{
		Id: id,
		Item: &pb.ToDoItem{
			Title:       todoItem.GetTitle(),
			Description: todoItem.GetDescription(),
			Status:      todoItem.GetStatus(),
		},
		CreatedAt: timestamppb.New(createdAt),
	}, nil
}

func main() {
	db, err := store.New()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterToDoServiceServer(s, &server{Store: db})

	reflection.Register(s)
	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
