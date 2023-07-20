package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	api "github.com/nayakunin/gophermart/internal/generated"
	"github.com/nayakunin/gophermart/internal/logger"
	"github.com/nayakunin/gophermart/internal/storage"
)

func (s Server) PostAPIUserRegister(w http.ResponseWriter, r *http.Request) *api.Response {
	response := api.Response{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return response.Status(http.StatusInternalServerError)
	}

	var req api.PostAPIUserRegisterJSONBody
	err = json.Unmarshal(body, &req)
	if err != nil {
		logger.Errorf("failed to unmarshal body: %v", err)
		return response.Status(http.StatusBadRequest)
	}

	if len(req.Login) == 0 || len(req.Password) == 0 {
		logger.Errorf("empty login or password")
		return response.Status(http.StatusBadRequest)
	}

	userID, err := s.Storage.CreateUser(req.Login, req.Password)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			logger.Errorf("failed to create user (user exists): %v", err)
			return response.Status(http.StatusConflict)
		}

		logger.Errorf("failed to create user: %v", err)
		return response.Status(http.StatusInternalServerError)
	}

	tokenString, err := s.TokenService.CreateToken(s.Cfg, userID)
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
