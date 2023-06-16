package server

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	api "github.com/nayakunin/gophermart/internal/generated"
	"github.com/nayakunin/gophermart/internal/logger"
	"github.com/nayakunin/gophermart/internal/storage"
	"github.com/nayakunin/gophermart/internal/utils/checksum"
)

func (s Server) PostAPIUserOrders(_ http.ResponseWriter, r *http.Request) *api.Response {
	response := api.Response{}
	userID := r.Context().Value(s.Cfg.AuthKey).(int64)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("failed to read body: %v", err)
		return response.Status(http.StatusInternalServerError)
	}

	orderID, err := strconv.ParseInt(string(body), 10, 64)
	if err != nil {
		logger.Errorf("failed to convert order to int: %v", err)
		return response.Status(http.StatusBadRequest)
	}

	if !checksum.Valid(orderID) {
		logger.Errorf("invalid checksum")
		return response.Status(http.StatusUnprocessableEntity)
	}

	err = s.Storage.SaveOrder(userID, orderID)
	if err != nil {
		if errors.Is(err, storage.ErrSaveOrderConflict) {
			logger.Errorf("failed to save order (conflict): %v", err)
			return response.Status(http.StatusConflict)
		}

		if errors.Is(err, storage.ErrSaveOrderAlreadyExists) {
			logger.Errorf("failed to save order (already exists): %v", err)
			return response.Status(http.StatusOK)
		}

		logger.Errorf("failed to save order: %v", err)
		return response.Status(http.StatusInternalServerError)
	}

	s.Worker.AddOrder(userID, orderID)

	return response.Status(http.StatusAccepted)
}
