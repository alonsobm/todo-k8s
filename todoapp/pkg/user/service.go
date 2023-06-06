package user

import (
	"context"
	"errors"
	"github.com/jackc/pgx"
	"time"
)

type User struct {
	UserID         int
	Email          string
	Username       string
	HashedPassword string
}

type Service interface {
	Create(ctx context.Context, cmd *CreateUserCommand) (*User, error)
	Delete(ctx context.Context, cmd *DeleteUserCommand) (*User, error)
	Login(ctx context.Context, cmd *LoginUserCommand) (*Auth, error)
}

type CreateUserCommand struct {
	UserName string
	Email    string
	Password string
}

type DeleteUserCommand struct {
	UserID int
}

type LoginUserCommand struct {
	UserEmail string
	Password  string
}

type ServiceImpl struct {
	conn *pgx.ConnPool
	auth *AuthService
}

func NewServiceImpl(conn *pgx.ConnPool, authSvc *AuthService) *ServiceImpl {
	return &ServiceImpl{
		auth: authSvc,
		conn: conn,
	}
}

func (s ServiceImpl) Create(ctx context.Context, cmd *CreateUserCommand) (*User, error) {

	var user User
	hashedPassword, err := s.auth.HashPassword(cmd.Password)
	if err != nil {
		return nil, err
	}

	query := `insert into users(username, email, hashed_password) values ($1,$2,$3) returning  user_id, username, email`
	row := s.conn.QueryRow(query, cmd.UserName, cmd.Email, hashedPassword)
	if row == nil {
		return nil, err
	}

	err = row.Scan(&user.UserID, &user.Username, &user.Email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s ServiceImpl) Delete(ctx context.Context, cmd *DeleteUserCommand) (*User, error) {

	query := `delete from users where user_id = $1 returning  user_id, email, username`

	row := s.conn.QueryRow(query, cmd.UserID)
	if row == nil {
		return nil, errors.New("error sql Delete")
	}

	user := new(User)
	err := row.Scan(&user.UserID, &user.Email, &user.Username)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return user, nil
}

func (s ServiceImpl) Login(ctx context.Context, cmd *LoginUserCommand) (*Auth, error) {

	query := `select user_id, email, username, hashed_password from users where email = $1 `

	row := s.conn.QueryRow(query, cmd.UserEmail)
	if row == nil {
		return nil, errors.New("no rows with this email")
	}

	var user User
	err := row.Scan(&user.UserID, &user.Email, &user.Username, &user.HashedPassword)
	if err != nil {
		return nil, err
	}

	err = s.auth.CheckPassword(cmd.Password, user.HashedPassword)
	if err != nil {
		return nil, err
	}

	token, err := s.auth.CreateToken(user.UserID, user.Username, time.Hour*48)
	if err != nil {
		return nil, err
	}

	return &Auth{
		token:     token,
		userID:    user.UserID,
		userEmail: user.Email,
	}, nil
}

type Auth struct {
	token     string
	userID    int
	userEmail string
}
