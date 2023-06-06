package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"kuberneteslab/todoapp/pkg/todo"
	"kuberneteslab/todoapp/pkg/user"
	"net/http"
	"testing"
)

func TestIntegrationTodos(t *testing.T) {

	t.Run("create a todo with valid token for myself should return ok", func(t *testing.T) {
		userID, token := credentialsHelper(t)
		defer func() {
			_, err := us.Delete(context.Background(), &user.DeleteUserCommand{UserID: userID})
			require.NoError(t, err)
		}()

		todoRequestDTO := todo.CreateTodoRequestDTO{
			UserID:  userID,
			Title:   "title1",
			Content: "content1",
		}
		marshalled, err := json.Marshal(&todoRequestDTO)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/todo", bytes.NewBuffer(marshalled))
		require.NoError(t, err)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		response, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		bytesReaded, err := io.ReadAll(response.Body)
		require.NoError(t, err)

		var todoResponseDTO todo.CreateTodoResponseDTO
		err = json.Unmarshal(bytesReaded, &todoResponseDTO)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, response.StatusCode)
		require.Equal(t, todoRequestDTO.Title, todoResponseDTO.Name)
		require.Equal(t, todoRequestDTO.Content, todoResponseDTO.Content)
		require.Greater(t, todoResponseDTO.TodoID, 0)

	})

	t.Run("create a todo with valid token for other user should return unauthorized", func(t *testing.T) {
		userID, token := credentialsHelper(t)
		defer func() {
			_, err := us.Delete(context.Background(), &user.DeleteUserCommand{UserID: userID})
			require.NoError(t, err)
		}()

		todoRequestDTO := todo.CreateTodoRequestDTO{
			UserID:  999,
			Title:   "title1",
			Content: "content1",
		}
		marshalled, err := json.Marshal(&todoRequestDTO)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/todo", bytes.NewBuffer(marshalled))
		require.NoError(t, err)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		response, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		bytesReaded, err := io.ReadAll(response.Body)
		require.NoError(t, err)

		t.Log(string(bytesReaded))
		require.Equal(t, http.StatusUnauthorized, response.StatusCode)
		require.Equal(t, "can't create todos for another user", string(bytesReaded))

	})

	t.Run("create a todo with invalid token should return unauthorized", func(t *testing.T) {
		userID, token := credentialsHelper(t)
		defer func() {
			_, err := us.Delete(context.Background(), &user.DeleteUserCommand{UserID: userID})
			require.NoError(t, err)
		}()
		token = "fake token"

		todoRequestDTO := todo.CreateTodoRequestDTO{
			UserID:  999,
			Title:   "title1",
			Content: "content1",
		}
		marshalled, err := json.Marshal(&todoRequestDTO)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/todo", bytes.NewBuffer(marshalled))
		require.NoError(t, err)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		response, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		bytesReaded, err := io.ReadAll(response.Body)
		require.NoError(t, err)

		t.Log(string(bytesReaded))
		require.Equal(t, http.StatusUnauthorized, response.StatusCode)
	})

	t.Run("update a todo with valid token that belongs to me should return ok", func(t *testing.T) {
		userID, token := credentialsHelper(t)
		defer func() {
			_, err := us.Delete(context.Background(), &user.DeleteUserCommand{UserID: userID})
			require.NoError(t, err)
		}()

		todoResponse := createTodoHelper(t, userID, token)
		defer func() {
			_, err := ts.Delete(context.Background(), &todo.DeleteTodoCommand{UserID: userID, TodoID: todoResponse.TodoID})
			require.NoError(t, err)
		}()

		updateRequest := todo.UpdateTodoRequestDTO{
			UserID:  userID,
			TodoID:  todoResponse.TodoID,
			Title:   "title updated",
			Content: "content updated",
		}
		marshalled, err := json.Marshal(&updateRequest)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPatch, "http://localhost:8080/todo", bytes.NewBuffer(marshalled))
		require.NoError(t, err)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		response, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		bytesReaded, err := io.ReadAll(response.Body)
		require.NoError(t, err)

		var todoResponseDTO todo.UpdateTodoResponseDTO
		err = json.Unmarshal(bytesReaded, &todoResponseDTO)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, response.StatusCode)
		require.Equal(t, updateRequest.Title, todoResponseDTO.Title)
		require.Equal(t, updateRequest.Content, todoResponseDTO.Content)
		require.Equal(t, updateRequest.TodoID, todoResponseDTO.TodoID)
		require.Equal(t, updateRequest.UserID, todoResponseDTO.UserID)
	})

	t.Run("update todo with valid token that belongs to other user should return unauthorized", func(t *testing.T) {
		userID, token := credentialsHelper(t)
		defer func() {
			_, err := us.Delete(context.Background(), &user.DeleteUserCommand{UserID: userID})
			require.NoError(t, err)
		}()

		todoResponse := createTodoHelper(t, userID, token)
		defer func() {
			_, err := ts.Delete(context.Background(), &todo.DeleteTodoCommand{UserID: userID, TodoID: todoResponse.TodoID})
			require.NoError(t, err)
		}()

		updateRequest := todo.UpdateTodoRequestDTO{
			UserID:  999,
			TodoID:  todoResponse.TodoID,
			Title:   "title updated",
			Content: "content updated",
		}
		marshalled, err := json.Marshal(&updateRequest)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPatch, "http://localhost:8080/todo", bytes.NewBuffer(marshalled))
		require.NoError(t, err)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		response, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusUnauthorized, response.StatusCode)

	})

	t.Run("update todo with invalid token should return unauthorized", func(t *testing.T) {
		userID, token := credentialsHelper(t)
		defer func() {
			_, err := us.Delete(context.Background(), &user.DeleteUserCommand{UserID: userID})
			require.NoError(t, err)
		}()

		todoResponse := createTodoHelper(t, userID, token)
		defer func() {
			_, err := ts.Delete(context.Background(), &todo.DeleteTodoCommand{UserID: userID, TodoID: todoResponse.TodoID})
			require.NoError(t, err)
		}()

		updateRequest := todo.UpdateTodoRequestDTO{
			UserID:  999,
			TodoID:  todoResponse.TodoID,
			Title:   "title updated",
			Content: "content updated",
		}
		marshalled, err := json.Marshal(&updateRequest)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPatch, "http://localhost:8080/todo", bytes.NewBuffer(marshalled))
		require.NoError(t, err)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "invalid token"))
		response, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusUnauthorized, response.StatusCode)

	})

	t.Run("get todo with valid token should return ok", func(t *testing.T) {
		userID, token := credentialsHelper(t)
		defer func() {
			_, err := us.Delete(context.Background(), &user.DeleteUserCommand{UserID: userID})
			require.NoError(t, err)
		}()

		todoResponse := createTodoHelper(t, userID, token)
		defer func() {
			_, err := ts.Delete(context.Background(), &todo.DeleteTodoCommand{UserID: userID, TodoID: todoResponse.TodoID})
			require.NoError(t, err)
		}()

		url := fmt.Sprintf("http://localhost:8080/todo?user_id=%v", userID)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		response, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		bytesReaded, err := io.ReadAll(response.Body)
		require.NoError(t, err)

		var todoResponseDTO todo.GetAllResponseDTO
		err = json.Unmarshal(bytesReaded, &todoResponseDTO)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, response.StatusCode)
		require.Equal(t, 1, len(todoResponseDTO.Todos))
		require.Equal(t, todoResponseDTO.Todos[0].TodoID, todoResponse.TodoID)
		require.Equal(t, todoResponseDTO.Todos[0].Title, todoResponse.Name)
		require.Equal(t, todoResponseDTO.Todos[0].Content, todoResponse.Content)

	})

	t.Run("get todo with valid token that belongs to other users should return unauthorized", func(t *testing.T) {
		userID, token := credentialsHelper(t)
		defer func() {
			_, err := us.Delete(context.Background(), &user.DeleteUserCommand{UserID: userID})
			require.NoError(t, err)
		}()

		todoResponse := createTodoHelper(t, userID, token)
		defer func() {
			_, err := ts.Delete(context.Background(), &todo.DeleteTodoCommand{UserID: userID, TodoID: todoResponse.TodoID})
			require.NoError(t, err)
		}()

		url := fmt.Sprintf("http://localhost:8080/todo?user_id=%v", 999)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		response, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusUnauthorized, response.StatusCode)
	})

	t.Run("get todo with invalid token should return unauthorized", func(t *testing.T) {
		userID, token := credentialsHelper(t)
		defer func() {
			_, err := us.Delete(context.Background(), &user.DeleteUserCommand{UserID: userID})
			require.NoError(t, err)
		}()

		createTodoResponse := createTodoHelper(t, userID, token)
		defer func() {
			_, err := ts.Delete(context.Background(), &todo.DeleteTodoCommand{UserID: userID, TodoID: createTodoResponse.TodoID})
			require.NoError(t, err)
		}()

		url := fmt.Sprintf("http://localhost:8080/todo?user_id=%v", userID)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "invalidtoken"))
		response, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusUnauthorized, response.StatusCode)
	})

	t.Run("delete todo with valid token that belongs to me should return ok", func(t *testing.T) {
		userID, token := credentialsHelper(t)
		defer func() {
			_, err := us.Delete(context.Background(), &user.DeleteUserCommand{UserID: userID})
			require.NoError(t, err)
		}()

		todoResponse := createTodoHelper(t, userID, token)

		deleteRequest := todo.DeleteTodoResquestDTO{
			UserID: userID,
			TodoID: todoResponse.TodoID,
		}
		marshalled, err := json.Marshal(&deleteRequest)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/todo", bytes.NewBuffer(marshalled))
		require.NoError(t, err)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		response, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		bytesReaded, err := io.ReadAll(response.Body)
		require.NoError(t, err)

		var todoResponseDTO todo.DeleteTodoResponseDTO
		err = json.Unmarshal(bytesReaded, &todoResponseDTO)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, response.StatusCode)
		require.Equal(t, todoResponse.TodoID, todoResponseDTO.TodoID)
		require.Equal(t, todoResponse.Content, todoResponseDTO.Content)
		require.Equal(t, todoResponse.Name, todoResponseDTO.Title)

	})

	t.Run("delete todo with valid token that belongs to other user should return unauthorized", func(t *testing.T) {
		userID, token := credentialsHelper(t)
		defer func() {
			_, err := us.Delete(context.Background(), &user.DeleteUserCommand{UserID: userID})
			require.NoError(t, err)
		}()

		todoResponse := createTodoHelper(t, userID, token)

		updateRequest := todo.DeleteTodoResquestDTO{
			UserID: 999,
			TodoID: todoResponse.TodoID,
		}
		marshalled, err := json.Marshal(&updateRequest)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/todo", bytes.NewBuffer(marshalled))
		require.NoError(t, err)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		response, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusUnauthorized, response.StatusCode)

	})

	t.Run("delete todo with invvalid token should return unauthorized", func(t *testing.T) {
		userID, token := credentialsHelper(t)
		defer func() {
			_, err := us.Delete(context.Background(), &user.DeleteUserCommand{UserID: userID})
			require.NoError(t, err)
		}()

		todoResponse := createTodoHelper(t, userID, token)

		updateRequest := todo.DeleteTodoResquestDTO{
			UserID: userID,
			TodoID: todoResponse.TodoID,
		}
		marshalled, err := json.Marshal(&updateRequest)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/todo", bytes.NewBuffer(marshalled))
		require.NoError(t, err)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "invalidtoken"))
		response, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusUnauthorized, response.StatusCode)

	})

}

func credentialsHelper(t *testing.T) (int, string) {

	_, err := us.Create(context.Background(), &user.CreateUserCommand{
		UserName: "username1",
		Email:    "email1",
		Password: "password1",
	})
	require.NoError(t, err)
	requestDTO := user.LoginUserRequestDTO{
		UserEmail:    "email1",
		UserPassword: "password1",
	}
	marshalled, err := json.Marshal(&requestDTO)
	require.NoError(t, err)

	res, err := http.Post("http://localhost:8080/user/login", "application/json", bytes.NewBuffer(marshalled))
	require.NoError(t, err)
	bytesReaded, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	var responseDTO user.LoginUserResponseDTO
	err = json.Unmarshal(bytesReaded, &responseDTO)
	require.NoError(t, err)
	return responseDTO.UserID, responseDTO.Token

}

func createTodoHelper(t *testing.T, userID int, token string) todo.CreateTodoResponseDTO {

	todoRequestDTO := todo.CreateTodoRequestDTO{
		UserID:  userID,
		Title:   "title1",
		Content: "content1",
	}
	marshalled, err := json.Marshal(&todoRequestDTO)
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/todo", bytes.NewBuffer(marshalled))
	require.NoError(t, err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	response, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	bytesReaded, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	var todoResponseDTO todo.CreateTodoResponseDTO
	err = json.Unmarshal(bytesReaded, &todoResponseDTO)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, response.StatusCode)
	require.Equal(t, todoRequestDTO.Title, todoResponseDTO.Name)
	require.Equal(t, todoRequestDTO.Content, todoResponseDTO.Content)
	require.Greater(t, todoResponseDTO.TodoID, 0)
	return todoResponseDTO
}
