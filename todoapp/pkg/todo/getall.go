package todo

import (
	"encoding/json"
	"kuberneteslab/todoapp/pkg/middlewares"
	"net/http"
	"strconv"
)

type GetAllTodoHttpHandler struct {
	Service Service
}

func NewGetAllTodoHttpHandler(s Service) *GetAllTodoHttpHandler {
	return &GetAllTodoHttpHandler{Service: s}
}

type GetAllResponseDTO struct {
	Todos []ResponseTodoDTO `json:"todos"`
}

type ResponseTodoDTO struct {
	TodoID  int    `json:"todo_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (h GetAllTodoHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
		return
	}

	cmd := GetAllTodosCommand{
		UserID: userID,
	}

	tokenUserID := r.Header.Get(middlewares.HeaderKeyUserID)
	strUserID := strconv.Itoa(userID)

	if tokenUserID != strUserID {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("can't get todos for another user"))
		return
	}

	todos, err := h.Service.GetAll(r.Context(), &cmd)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	res := make([]ResponseTodoDTO, 0, len(todos))
	for _, todo := range todos {
		res = append(res, ResponseTodoDTO{
			TodoID:  todo.TodoID,
			Title:   todo.Title,
			Content: todo.Content,
		})
	}

	responseDTO := GetAllResponseDTO{
		Todos: res,
	}
	bytes, err := json.Marshal(responseDTO)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
