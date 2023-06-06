package middlewares

import (
	"kuberneteslab/todoapp/pkg/user"
	"net/http"
	"strconv"
	"strings"
)

const (
	HeaderKeyUserID   = "userID"
	HeaderKeyUserName = "userName"
)

func AuthMiddleware(auth *user.AuthService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		header := r.Header.Get("Authorization")
		if len(header) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("auth token not provided"))
			return
		}

		fields := strings.Fields(header)
		if len(fields) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("auth token bad format"))
			return
		}

		if fields[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("auth token not bearer"))
			return
		}

		payload, err := auth.VerifyToken(fields[1])
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("auth token couldn't be verified"))
			return
		}

		r.Header.Set(HeaderKeyUserID, strconv.Itoa(payload.UserID))
		r.Header.Set(HeaderKeyUserName, payload.Username)

		next.ServeHTTP(w, r)
	})
}
