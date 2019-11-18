package v1

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	v1 "github.com/sauravgsh16/api-grpc/pkg/api/v1"
)

const (
	apiversion = "v1"
)

type toDoServiceServer struct {
	db *sql.DB
}

// NewToDoServiceServer returns a new toDoSeriveServer
func NewToDoServiceServer(db *sql.DB) v1.ToDoServiceServer {
	return &toDoServiceServer{
		db: db,
	}
}

func (s *toDoServiceServer) checkAPI(api string) error {
	if len(api) > 0 {
		if apiversion != api {
			return status.Error(codes.Unimplemented, "unsupported API version")
		}
	}
	return nil
}

func (s *toDoServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("failed to connect to database: %v", err))
	}
	return c, nil
}

func (s *toDoServiceServer) Create(ctx context.Context, req *v1.CreateRequest) (*v1.CreateResponse, error) {
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	timestamp, err := ptypes.Timestamp(req.ToDo.Reminder)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid format for reminder: "+err.Error())
	}

	var returnedID int64

	err = c.QueryRowContext(ctx, `INSERT INTO ToDo (Title, Description, Reminder) VALUES($1, $2, $3) RETURNING ID`,
		req.ToDo.Title, req.ToDo.Description, timestamp).Scan(&returnedID)
	if err != nil {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("failed to insert row: %v", err))
	}

	return &v1.CreateResponse{
		Api: apiversion,
		Id:  returnedID,
	}, nil
}

func (s *toDoServiceServer) Read(ctx context.Context, req *v1.ReadRequest) (*v1.ReadResponse, error) {
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	rows, err := c.QueryContext(ctx, `SELECT ID, Title, Description, Reminder FROM ToDO WHERE ID=$1`, req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from ToDo: "+err.Error())
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve data from ToDo: "+err.Error())
		}
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ToDo with ID='%d' is not found", req.Id))
	}

	// Get data from row and create response
	var td v1.ToDo
	var reminder time.Time

	if err := rows.Scan(&td.Id, &td.Title, &td.Description, &reminder); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve value from rows: "+err.Error())
	}
	td.Reminder, err = ptypes.TimestampProto(reminder)
	if err != nil {
		return nil, status.Error(codes.Unknown, "invalid format for reminder field: "+err.Error())
	}

	// Check if any other rows exists. Error if present
	if rows.Next() {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("found multiple rows in ToDo with ID='%d'", req.Id))
	}

	return &v1.ReadResponse{
		Api:  apiversion,
		Todo: &td,
	}, nil
}

func (s *toDoServiceServer) Update(ctx context.Context, req *v1.UpdateRequest) (*v1.UpdateResponse, error) {
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	reminder, err := ptypes.Timestamp(req.Todo.Reminder)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid format for reminder field: "+err.Error())
	}

	// Execute update query
	resp, err := c.ExecContext(ctx, `UPDATE ToDo Set Title=$1, Description=$2, Reminder=$3 WHERE ID=$4`,
		req.Todo.Title, req.Todo.Description, reminder, req.Todo.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to update ToDo: "+err.Error())
	}

	rows, err := resp.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected: "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Task with ID='%d', not found", req.Todo.Id))
	}

	return &v1.UpdateResponse{
		Api:     apiversion,
		Updated: rows,
	}, nil
}

func (s *toDoServiceServer) Delete(ctx context.Context, req *v1.DeleteRequest) (*v1.DeleteResponse, error) {
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	resp, err := c.ExecContext(ctx, `DELETE FROM ToDo WHERE ID=$1`, req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to delete row: "+err.Error())
	}

	rows, err := resp.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected: "+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Row with ID='%d', not found", req.Id))
	}

	return &v1.DeleteResponse{
		Api:     apiversion,
		Deleted: rows,
	}, nil
}

func (s *toDoServiceServer) ReadAll(ctx context.Context, req *v1.ReadAllRequest) (*v1.ReadAllResponse, error) {
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	rows, err := c.QueryContext(ctx, `SELECT Id, Title, Description, Reminder FROM ToDo`)
	if err != nil {
		return nil, status.Error(codes.Unknown, "fail to select rows from ToDo: "+err.Error())
	}
	defer rows.Close()

	var reminder time.Time
	todoList := []*v1.ToDo{}

	for rows.Next() {
		td := new(v1.ToDo)
		if err := rows.Scan(&td.Id, &td.Title, &td.Description, &reminder); err != nil {
			return nil, status.Error(codes.Unknown, "failed to scan rows: "+err.Error())
		}
		td.Reminder, err = ptypes.TimestampProto(reminder)
		if err != nil {
			return nil, status.Error(codes.Unknown, "reminder field -invalid format: "+err.Error())
		}
		todoList = append(todoList, td)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve data from rows: "+err.Error())
	}

	resp := &v1.ReadAllResponse{
		Api:   apiversion,
		ToDos: todoList,
	}

	return resp, nil
}
