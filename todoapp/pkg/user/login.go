package user

import (
	"encoding/json"
	"io"
	"net/http"
)

type LoginUserHttpHandler struct {
	Service Service
}

func NewLoginUserHttpHandler(s Service) *LoginUserHttpHandler {
	return &LoginUserHttpHandler{Service: s}
}

type LoginUserRequestDTO struct {
	UserEmail    string `json:"email"`
	UserPassword string `json:"password"`
}

type LoginUserResponseDTO struct {
	Token  string `json:"token"`
	UserID int    `json:"user_id"`
}

func (h LoginUserHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var request LoginUserRequestDTO
	err = json.Unmarshal(bytes, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cmd := LoginUserCommand{
		UserEmail: request.UserEmail,
		Password:  request.UserPassword,
	}

	auth, err := h.Service.Login(r.Context(), &cmd)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	response := LoginUserResponseDTO{Token: auth.token, UserID: auth.userID}
	bytes, err = json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
	return
}
