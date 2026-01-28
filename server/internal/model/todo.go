package model

import "time"

type ToDoItem struct {
	ID          int64
	Title       string
	Description string
	Status      int32
	CreatedAt   time.Time
}
