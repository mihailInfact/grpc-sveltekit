package main

import (
	"context"
	"greeter/internal/store"
	pb "greeter/pkg/greeter"
	"greeter/pkg/greeter/greeterconnect"
	"log"
	"net/http"
	"strings"
	"time"

	"connectrpc.com/connect"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:generate sh -c "cd .. && buf generate --path proto/greeter.proto"

type server struct {
	pb.UnimplementedToDoServiceServer
	Store *store.Store
}

func (s *server) GetOne(ctx context.Context, req *connect.Request[pb.GetOneRequest]) (*connect.Response[pb.GetOneResponse], error) {
	return connect.NewResponse(&pb.GetOneResponse{}), nil
}

func (s *server) Update(ctx context.Context, req *connect.Request[pb.UpdateRequest]) (*connect.Response[pb.UpdateResponse], error) {
	return connect.NewResponse(&pb.UpdateResponse{}), nil
}

func (s *server) GetAll(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[pb.GetAllResponse], error) {
	log.Println("fetching all todo items")

	connectVersion := req.Header().Get("Connect-Protocol-Version")

	// 2. Check the Content-Type to distinguish between gRPC and gRPC-Web
	contentType := req.Header().Get("Content-Type")

	if connectVersion != "" {
		log.Printf("Protocol Used: Connect (Version %s)", connectVersion)
	} else if contentType == "application/grpc" {
		log.Println("Protocol Used: Standard gRPC")
	} else if contentType == "application/grpc-web" {
		log.Println("Protocol Used: gRPC-Web")
	} else {
		log.Printf("Protocol Used: Unknown (Content-Type: %s)", contentType)
	}

	if strings.HasPrefix(contentType, "application/grpc") {
		log.Println("Protocol Used: Standard gRPC (Binary)")
	} else if strings.Contains(contentType, "application/connect") {
		log.Println("Protocol Used: Connect Protocol")
	}

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

	return connect.NewResponse(&pb.GetAllResponse{Items: items}), nil
}

func (s *server) Create(ctx context.Context, req *connect.Request[pb.CreateRequest]) (*connect.Response[pb.CreateResponse], error) {
	todoItem := req.Msg.GetItem()
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
	return connect.NewResponse(&pb.CreateResponse{Item: &pb.ToDoItem{Id: id,
		Item: &pb.ToDoDetails{
			Title:       todoItem.GetTitle(),
			Description: todoItem.GetDescription(),
			Status:      todoItem.GetStatus(),
		},
		CreatedAt: timestamppb.New(createdAt)}}), nil
}

func (s *server) Delete(ctx context.Context, req *connect.Request[pb.DeleteRequest]) (*connect.Response[emptypb.Empty], error) {
	id := req.Msg.GetId()
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

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func (s *server) UpdateStatus(ctx context.Context, req *connect.Request[pb.UpdateStatusRequest]) (*connect.Response[emptypb.Empty], error) {
	id := req.Msg.Id
	newStatus := req.Msg.Status

	log.Printf("Updating todo %v to status %v", id, newStatus)

	_, err := s.Store.ExecContext(ctx, "UPDATE todos SET status = ? WHERE id = ?", newStatus, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not update: %v", err)
	}

	return connect.NewResponse(&emptypb.Empty{}), nil
}

func main() {
	db, err := store.New()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	srv := &server{Store: db}
	mux := http.NewServeMux()
	path, handler := greeterconnect.NewToDoServiceHandler(srv)
	mux.Handle(path, handler)

	log.Printf("Connect server listening at :50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	err = http.ListenAndServe(
		"localhost:50051",
		h2c.NewHandler(mux, &http2.Server{}),
	)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
