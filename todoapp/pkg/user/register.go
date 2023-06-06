package user

import (
	"encoding/json"
	"io"
	"net/http"
)

type CreateUserHttpHandler struct {
	Service Service
}

func NewCreateUserHttpHandler(s Service) *CreateUserHttpHandler {
	return &CreateUserHttpHandler{Service: s}
}

type CreateUserRequestDTO struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResponseDTO struct {
	UserID   int    `json:"user_id"`
	UserName string `json:"title"`
	Email    string `json:"email"`
}

func (h CreateUserHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	dto := CreateUserRequestDTO{}
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

	cmd := CreateUserCommand{
		Email:    dto.Email,
		Password: dto.Password,
		UserName: dto.UserName,
	}

	user, err := h.Service.Create(r.Context(), &cmd)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bytes, err = json.Marshal(CreateUserResponseDTO{
		UserID:   user.UserID,
		UserName: user.Username,
		Email:    user.Email,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
	return
}
