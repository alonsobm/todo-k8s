package server

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx"
	"kuberneteslab/todoapp/pkg/middlewares"
	"kuberneteslab/todoapp/pkg/todo"
	"kuberneteslab/todoapp/pkg/user"
	"log"
	"net/http"
	"os"
	"time"
)

func Start(configPath string) {
	var config Config
	content, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal("error reading file ", err.Error())
	}

	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatal("invalid config file: ", err.Error())
	}

	var tlscfg *tls.Config
	if config.Environment == "prod" {
		tlscfg = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	pc := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:      config.DatabaseHost,
			Port:      uint16(config.DatabasePort),
			Database:  config.DatabaseName,
			User:      config.DatabaseUser,
			Password:  config.DatabasePassword,
			TLSConfig: tlscfg,
		},
	}

	conn, err := pgx.NewConnPool(pc)
	if err != nil {
		log.Fatal("cannot connect to database: ", err.Error())
	}

	r := chi.NewRouter()

	authSvc, err := user.NewAuthService(config.AuthKey)
	if err != nil {
		log.Fatal("error creating auth service", err.Error())
	}
	us := user.NewServiceImpl(conn, authSvc)
	ts := todo.NewServiceImpl(conn)

	createUser := user.NewCreateUserHttpHandler(us)
	loginUser := user.NewLoginUserHttpHandler(us)

	createTodo := todo.NewCreateTodoHttpHandler(ts)
	deleteTodo := todo.NewDeleteTodoHttpHandler(ts)
	updateTodo := todo.NewUpdateTodoHttpHandler(ts)
	getAllTodo := todo.NewGetAllTodoHttpHandler(ts)

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Get("/", Hello)
	r.Get("/ping", Ping)

	r.Post("/user", createUser.ServeHTTP)
	r.Post("/user/login", loginUser.ServeHTTP)
	r.Post("/todo", middlewares.AuthMiddleware(authSvc, createTodo).ServeHTTP)
	r.Delete("/todo", middlewares.AuthMiddleware(authSvc, deleteTodo).ServeHTTP)
	r.Get("/todo", middlewares.AuthMiddleware(authSvc, getAllTodo).ServeHTTP)
	r.Patch("/todo", middlewares.AuthMiddleware(authSvc, updateTodo).ServeHTTP)

	log.Println("lets listen")
	err = http.ListenAndServe(config.Port, r)
	if err != nil {
		log.Println("error listening: ", err.Error())
	}
}

func Hello(w http.ResponseWriter, r *http.Request) {
	name, _ := os.Hostname()
	template := "Hello Kubernetes. Time: %v. Pod: %v.\n"
	response := []byte(fmt.Sprintf(template, time.Now(), name))
	w.Write(response)
}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
