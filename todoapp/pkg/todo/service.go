package todo

import (
	"context"
	"errors"
	"github.com/jackc/pgx"
	"time"
)

type Todo struct {
	TodoID  int
	UserID  int
	Title   string
	Content string
}

type Service interface {
	GetAll(ctx context.Context, cmd *GetAllTodosCommand) ([]*Todo, error)
	Create(ctx context.Context, cmd *CreateTodoCommand) (*Todo, error)
	Update(ctx context.Context, cmd *UpdateTodoCommand) (*Todo, error)
	Get(ctx context.Context, cmd *GetTodoCommand) (*Todo, error)
	Delete(ctx context.Context, cmd *DeleteTodoCommand) (*Todo, error)
}

type GetAllTodosCommand struct {
	UserID int
}

type GetTodoCommand struct {
	UserID int
	TodoID int
}

type CreateTodoCommand struct {
	UserID  int
	Title   string
	Content string
}

type UpdateTodoCommand struct {
	UserID  int
	TodoID  int
	title   string
	content string
}
type DeleteTodoCommand struct {
	UserID int
	TodoID int
}

type ServiceImpl struct {
	conn *pgx.ConnPool
}

func NewServiceImpl(conn *pgx.ConnPool) *ServiceImpl {
	return &ServiceImpl{
		conn: conn,
	}
}

func (s ServiceImpl) Create(ctx context.Context, cmd *CreateTodoCommand) (*Todo, error) {

	todo := Todo{}
	query := `insert into todos(user_id, title, content) values ($1,$2,$3) returning todo_id, user_id, title, content`
	row := s.conn.QueryRow(query, cmd.UserID, cmd.Title, cmd.Content)
	if row == nil {
		return nil, errors.New("err todo create empty row")
	}
	err := row.Scan(&todo.TodoID, &todo.UserID, &todo.Title, &todo.Content)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (s ServiceImpl) GetAll(ctx context.Context, cmd *GetAllTodosCommand) ([]*Todo, error) {

	query := `select todo_id, title, content from todos where user_id = $1`

	rows, err := s.conn.Query(query, cmd.UserID)
	if rows == nil {
		return nil, errors.New("error todo get all empty rows")
	}
	if err != nil {
		return nil, err
	}

	todos := make([]*Todo, 0)

	for rows.Next() {
		todo := Todo{}
		rows.Scan(&todo.TodoID, &todo.Title, &todo.Content)
		todos = append(todos, &todo)
	}

	return todos, nil
}

func (s ServiceImpl) Get(ctx context.Context, cmd *GetTodoCommand) (*Todo, error) {
	query := `select todo_id, title, content from todos where user_id = $1 and todo_id = $2`

	row := s.conn.QueryRow(query, cmd.UserID, cmd.TodoID)
	if row == nil {
		return nil, errors.New("error GetOne")
	}

	var todo Todo
	err := row.Scan(&todo.TodoID, &todo.UserID, &todo.Title, todo.Content)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (s ServiceImpl) Update(ctx context.Context, cmd *UpdateTodoCommand) (*Todo, error) {

	query := `update todos set title = $1, content = $2, updated_at = $3 where todo_id = $4 returning  todo_id, user_id, title, content`
	row := s.conn.QueryRow(query, cmd.title, cmd.content, time.Now(), cmd.TodoID)
	if row == nil {
		return nil, errors.New("err todo update empty row")
	}
	var todo Todo
	err := row.Scan(&todo.TodoID, &todo.UserID, &todo.Title, &todo.Content)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (s ServiceImpl) Delete(ctx context.Context, cmd *DeleteTodoCommand) (*Todo, error) {

	query := `delete from todos where todo_id = $1 returning  todo_id, user_id, title, content`

	row := s.conn.QueryRow(query, cmd.TodoID)
	if row == nil {
		return nil, errors.New("error sql delete empty row")
	}

	todo := Todo{}
	err := row.Scan(&todo.TodoID, &todo.UserID, &todo.Title, &todo.Content)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}
