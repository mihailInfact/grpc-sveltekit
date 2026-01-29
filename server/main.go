package main

import (
	"context"
	"greeter/internal/store"
	pb "greeter/pkg/greeter"
	"log"
	"net"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:generate sh -c "cd .. && buf generate --path proto/greeter.proto"

type server struct {
	pb.UnimplementedToDoServiceServer
	Store *store.Store
}

func (s *server) GetAll(ctx context.Context, _ *emptypb.Empty) (*pb.GetAllResponse, error) {
	log.Println("fetching all todo items")

	items := make([]*pb.ToDoItem, 0)

	rows, err := s.Store.QueryContext(ctx,
		"SELECT id, title, description, status, created_at FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var title, description string
		var status int32
		var createdAt time.Time

		err := rows.Scan(&id, &title, &description, &status, &createdAt)
		if err != nil {
			log.Printf("Failed to scan row: %v", err)
			continue
		}

		item := &pb.ToDoItem{
			Id: id,
			Item: &pb.ToDoDetails{
				Title:       title,
				Description: description,
				Status:      pb.Status(status),
			},
			CreatedAt: timestamppb.New(createdAt),
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &pb.GetAllResponse{
		Items: items,
	}, nil
}

func (s *server) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	todoItem := req.GetItem()

	log.Printf("creating item: %+v\n", todoItem)

	var id int64
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
		Item: &pb.ToDoItem{
			Id: id,
			Item: &pb.ToDoDetails{
				Title:       todoItem.GetTitle(),
				Description: todoItem.GetDescription(),
				Status:      todoItem.GetStatus(),
			},
			CreatedAt: timestamppb.New(createdAt),
		},
	}, nil
}

func (s *server) Delete(ctx context.Context, req *pb.DeleteRequest) (*emptypb.Empty, error) {
	id := req.GetId()
	log.Println("deleting todo item: ", id)

	res, err := s.Store.ExecContext(ctx, "DELETE FROM todos WHERE ID = ?", id)
	if err != nil {
		return nil, err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		// Return a specific gRPC error code
		return nil, status.Errorf(codes.NotFound, "Todo with ID %d not found", id)
	}

	return &emptypb.Empty{}, nil
}

func (s *server) UpdateStatus(ctx context.Context, req *pb.UpdateStatusRequest) (*emptypb.Empty, error) {
	id := req.GetId()
	newStatus := req.GetStatus()

	log.Printf("Updating todo %v to status %v", id, newStatus)

	_, err := s.Store.ExecContext(ctx, "UPDATE todos SET status = ? WHERE id = ?", newStatus, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not update: %v", err)
	}

	return &emptypb.Empty{}, nil
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
