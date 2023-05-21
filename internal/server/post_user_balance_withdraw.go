package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	api "github.com/nayakunin/gophermart/internal/generated"
	"github.com/nayakunin/gophermart/internal/middlewares"
	"github.com/nayakunin/gophermart/internal/storage"
)

func (s Server) PostAPIUserBalanceWithdraw(_ http.ResponseWriter, r *http.Request) *api.Response {
	response := api.Response{}

	userID := r.Context().Value(middlewares.AuthKey).(int64)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return response.Status(http.StatusInternalServerError)
	}

	var req api.PostAPIUserBalanceWithdrawJSONRequestBody
	err = json.Unmarshal(body, &req)
	if err != nil {
		return response.Status(http.StatusBadRequest)
	}

	if len(req.Order) == 0 || req.Sum == 0 {
		return response.Status(http.StatusBadRequest)
	}

	orderID, err := strconv.Atoi(req.Order)
	if err != nil {
		return response.Status(http.StatusBadRequest)
	}

	err = s.Storage.Withdraw(userID, int64(orderID), req.Sum)
	if err != nil {
		if errors.Is(err, storage.ErrWithdrawOrderNotFound) {
			return response.Status(http.StatusUnprocessableEntity)
		}

		if errors.Is(err, storage.ErrWithdrawBalanceNotEnough) {
			return response.Status(http.StatusPaymentRequired)
		}

		return response.Status(http.StatusInternalServerError)
	}

	return response.Status(http.StatusOK)
}
