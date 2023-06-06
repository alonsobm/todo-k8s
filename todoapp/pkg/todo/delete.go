package todo

import (
	"encoding/json"
	"io"
	"kuberneteslab/todoapp/pkg/middlewares"
	"net/http"
	"strconv"
)

type DeleteTodoHttpHandler struct {
	Service Service
}

func NewDeleteTodoHttpHandler(s Service) *DeleteTodoHttpHandler {
	return &DeleteTodoHttpHandler{Service: s}
}

type DeleteTodoResquestDTO struct {
	UserID int `json:"user_id"`
	TodoID int `json:"todo_id"`
}

type DeleteTodoResponseDTO struct {
	TodoID  int    `json:"todo_id"`
	UserID  int    `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (h DeleteTodoHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	dto := DeleteTodoResquestDTO{}
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
		w.Write([]byte("can't delete todos for another user"))
		return
	}

	cmd := DeleteTodoCommand{
		TodoID: dto.TodoID,
		UserID: dto.UserID,
	}

	todo, err := h.Service.Delete(r.Context(), &cmd)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	responseDTO := DeleteTodoResponseDTO{
		TodoID:  todo.TodoID,
		UserID:  todo.UserID,
		Title:   todo.Title,
		Content: todo.Content,
	}

	bytes, err = json.Marshal(responseDTO)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
	return
}
