package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/nayakunin/gophermart/internal/auth"
	api "github.com/nayakunin/gophermart/internal/generated"
	"github.com/nayakunin/gophermart/internal/storage"
)

func (s Server) PostAPIUserLogin(w http.ResponseWriter, r *http.Request) *api.Response {
	response := api.Response{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return response.Status(http.StatusInternalServerError)
	}

	var req api.PostAPIUserLoginJSONBody
	err = json.Unmarshal(body, &req)
	if err != nil {
		return response.Status(http.StatusBadRequest)
	}

	if len(req.Login) == 0 || len(req.Password) == 0 {
		return response.Status(http.StatusBadRequest)
	}

	userID, err := s.Storage.GetUserID(req.Login, req.Password)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return response.Status(http.StatusUnauthorized)
		}
		return response.Status(http.StatusInternalServerError)
	}

	tokenString, err := auth.CreateToken(userID, s.Cfg.JWTSecret)
	if err != nil {
		return response.Status(http.StatusInternalServerError)
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "Authentication",
		Value: tokenString,
	})

	return response.Status(http.StatusOK)
}
