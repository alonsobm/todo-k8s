package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"io"
	"kuberneteslab/todoapp/pkg/user"
	"net/http"
	"testing"
)

func TestIntegrationUsers(t *testing.T) {

	t.Run("create user with a non used email or username should return ok", func(t *testing.T) {

		requestDTO := user.CreateUserRequestDTO{
			UserName: "username",
			Email:    "email",
			Password: "password",
		}
		marshalled, err := json.Marshal(&requestDTO)
		require.NoError(t, err)

		res, err := http.Post("http://localhost:8080/user", "application/json", bytes.NewBuffer(marshalled))
		require.NoError(t, err)
		bytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		var responseDTO user.CreateUserResponseDTO
		json.Unmarshal(bytes, &responseDTO)

		require.Equal(t, http.StatusOK, res.StatusCode)
		require.Equal(t, requestDTO.UserName, responseDTO.UserName)
		require.Greater(t, responseDTO.UserID, 0)
		require.Equal(t, requestDTO.Email, responseDTO.Email)

		_, err = us.Delete(context.Background(), &user.DeleteUserCommand{UserID: responseDTO.UserID})
		require.NoError(t, err)

	})

	t.Run("create user with an used email or username should return bad request", func(t *testing.T) {
		u, err := us.Create(context.Background(), &user.CreateUserCommand{
			UserName: "username",
			Email:    "email",
			Password: "password",
		})
		require.NoError(t, err)
		defer func() {
			_, err = us.Delete(context.Background(), &user.DeleteUserCommand{UserID: u.UserID})
			require.NoError(t, err)
		}()

		requestDTO := user.CreateUserRequestDTO{
			UserName: "username",
			Email:    "email",
			Password: "password",
		}
		marshalled, err := json.Marshal(&requestDTO)
		require.NoError(t, err)

		res, err := http.Post("http://localhost:8080/user", "application/json", bytes.NewBuffer(marshalled))
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, res.StatusCode)

	})

	t.Run("login user with a good password should return ok", func(t *testing.T) {

		u, err := us.Create(context.Background(), &user.CreateUserCommand{
			UserName: "username",
			Email:    "email",
			Password: "password",
		})
		require.NoError(t, err)
		defer func() {
			_, err = us.Delete(context.Background(), &user.DeleteUserCommand{UserID: u.UserID})
			require.NoError(t, err)
		}()

		requestDTO := user.LoginUserRequestDTO{
			UserEmail:    "email",
			UserPassword: "password",
		}
		marshalled, err := json.Marshal(&requestDTO)
		require.NoError(t, err)

		res, err := http.Post("http://localhost:8080/user/login", "application/json", bytes.NewBuffer(marshalled))
		require.NoError(t, err)
		bytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		var responseDTO user.LoginUserResponseDTO
		json.Unmarshal(bytes, &responseDTO)

		require.Equal(t, http.StatusOK, res.StatusCode)
		require.Equal(t, u.UserID, responseDTO.UserID)
		require.Greater(t, len(responseDTO.Token), 0)

	})

	t.Run("login user with an invalid password should return unauthorized", func(t *testing.T) {

		u, err := us.Create(context.Background(), &user.CreateUserCommand{
			UserName: "username",
			Email:    "email",
			Password: "password",
		})
		require.NoError(t, err)
		defer func() {
			_, err = us.Delete(context.Background(), &user.DeleteUserCommand{UserID: u.UserID})
			require.NoError(t, err)
		}()

		requestDTO := user.LoginUserRequestDTO{
			UserEmail:    "email",
			UserPassword: "invalid_password",
		}
		marshalled, err := json.Marshal(&requestDTO)
		require.NoError(t, err)

		res, err := http.Post("http://localhost:8080/user/login", "application/json", bytes.NewBuffer(marshalled))
		require.NoError(t, err)
		bytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		var responseDTO user.LoginUserResponseDTO
		json.Unmarshal(bytes, &responseDTO)

		require.Equal(t, http.StatusUnauthorized, res.StatusCode)

	})

}
