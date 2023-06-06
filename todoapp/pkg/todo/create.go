package todo

import (
	"encoding/json"
	"fmt"
	"io"
	"kuberneteslab/todoapp/pkg/middlewares"
	"net/http"
	"strconv"
)

type CreateTodoHttpHandler struct {
	Service Service
}

func NewCreateTodoHttpHandler(s Service) *CreateTodoHttpHandler {
	return &CreateTodoHttpHandler{Service: s}
}

type CreateTodoRequestDTO struct {
	UserID  int    `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type CreateTodoResponseDTO struct {
	TodoID  int    `json:"todo_id"`
	Name    string `json:"title"`
	Content string `json:"content"`
}

func (h CreateTodoHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	dto := CreateTodoRequestDTO{}
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(bytes, &dto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	tokenUserID := r.Header.Get(middlewares.HeaderKeyUserID)
	userID := strconv.Itoa(dto.UserID)

	if tokenUserID != userID {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("can't create todos for another user"))
		return
	}

	cmd := CreateTodoCommand{
		UserID:  dto.UserID,
		Title:   dto.Title,
		Content: dto.Content,
	}

	todo, err := h.Service.Create(r.Context(), &cmd)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error todo service create %v", err.Error())))
		return
	}

	responseDTO := CreateTodoResponseDTO{
		TodoID:  todo.TodoID,
		Name:    todo.Title,
		Content: todo.Content,
	}
	bytes, err = json.Marshal(responseDTO)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
