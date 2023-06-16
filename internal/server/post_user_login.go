package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/nayakunin/gophermart/internal/auth"
	api "github.com/nayakunin/gophermart/internal/generated"
	"github.com/nayakunin/gophermart/internal/logger"
	"github.com/nayakunin/gophermart/internal/storage"
)

func (s Server) PostAPIUserLogin(w http.ResponseWriter, r *http.Request) *api.Response {
	response := api.Response{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("failed to read request body: %v", err)
		return response.Status(http.StatusInternalServerError)
	}

	var req api.PostAPIUserLoginJSONBody
	err = json.Unmarshal(body, &req)
	if err != nil {
		logger.Errorf("failed to unmarshal request body: %v", err)
		return response.Status(http.StatusBadRequest)
	}

	if len(req.Login) == 0 || len(req.Password) == 0 {
		logger.Errorf("empty login or password")
		return response.Status(http.StatusBadRequest)
	}

	userID, err := s.Storage.GetUserID(req.Login, req.Password)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			logger.Errorf("user not found: %v", err)
			return response.Status(http.StatusUnauthorized)
		}

		logger.Errorf("failed to get user id: %v", err)
		return response.Status(http.StatusInternalServerError)
	}

	tokenString, err := auth.CreateToken(s.Cfg, userID)
	if err != nil {
		logger.Errorf("failed to create token: %v", err)
		return response.Status(http.StatusInternalServerError)
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "Authentication",
		Value: tokenString,
	})

	return response.Status(http.StatusOK)
}
