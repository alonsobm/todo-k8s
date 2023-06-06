package tests

import (
	"github.com/jackc/pgx"
	"kuberneteslab/todoapp/pkg/server"
	"kuberneteslab/todoapp/pkg/todo"
	"kuberneteslab/todoapp/pkg/user"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

var (
	us *user.ServiceImpl
	ts *todo.ServiceImpl
)

func TestMain(m *testing.M) {
	go func() {
		server.Start("../../config/local.json")
	}()

	for i := 0; i < 3; i++ {
		time.Sleep(time.Second * 1)
		response, err := http.Get("http://localhost:8080/ping")
		if err != nil {
			continue
		}

		if response.StatusCode == http.StatusOK {
			break
		}
	}

	pc := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:      "localhost",
			Port:      5432,
			Database:  "postgres",
			User:      "postgres",
			Password:  "postgres",
			TLSConfig: nil,
		},
	}

	conn, err := pgx.NewConnPool(pc)
	if err != nil {
		log.Fatal("cannot connect to database: ", err.Error())
	}
	authSvc, err := user.NewAuthService("12345678901234567890123456789012")
	if err != nil {
		log.Fatal("error creating auth service", err.Error())
	}

	us = user.NewServiceImpl(conn, authSvc)
	ts = todo.NewServiceImpl(conn)
	exitVal := m.Run()
	os.Exit(exitVal)
}
