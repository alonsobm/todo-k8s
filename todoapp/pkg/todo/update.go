package todo

import (
	"encoding/json"
	"io"
	"kuberneteslab/todoapp/pkg/middlewares"
	"net/http"
	"strconv"
)

type UpdateTodoHttpHandler struct {
	Service Service
}

func NewUpdateTodoHttpHandler(s Service) *UpdateTodoHttpHandler {
	return &UpdateTodoHttpHandler{Service: s}
}

type UpdateTodoRequestDTO struct {
	UserID  int    `json:"user_id"`
	TodoID  int    `json:"todo_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdateTodoResponseDTO struct {
	TodoID  int    `json:"todo_id"`
	UserID  int    `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (h UpdateTodoHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	dto := UpdateTodoRequestDTO{}
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
		w.Write([]byte("can't update todos for another user"))
		return
	}

	cmd := UpdateTodoCommand{
		UserID:  dto.UserID,
		TodoID:  dto.TodoID,
		title:   dto.Title,
		content: dto.Content,
	}

	todo, err := h.Service.Update(r.Context(), &cmd)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	responseDTO := UpdateTodoResponseDTO{
		TodoID:  todo.TodoID,
		UserID:  todo.UserID,
		Title:   todo.Title,
		Content: todo.Content,
	}

	bytes, err = json.Marshal(responseDTO)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
