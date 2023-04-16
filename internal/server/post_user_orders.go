package server

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	api "github.com/nayakunin/gophermart/internal/generated"
	"github.com/nayakunin/gophermart/internal/storage"
	"github.com/nayakunin/gophermart/internal/utils/checksum"
)

func (s Server) PostAPIUserOrders(_ http.ResponseWriter, r *http.Request) *api.Response {
	response := api.Response{}
	userID := r.Context().Value("userID").(int64)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return response.Status(http.StatusInternalServerError)
	}

	orderID, err := strconv.ParseInt(string(body), 10, 64)
	if err != nil {
		return response.Status(http.StatusBadRequest)
	}

	if !checksum.Valid(orderID) {
		return response.Status(http.StatusUnprocessableEntity)
	}

	err = s.Storage.SaveOrder(userID, orderID, api.OrderStatusNEW.ToValue())
	if err != nil {
		if errors.Is(err, storage.ErrSaveOrderConflict) {
			return response.Status(http.StatusConflict)
		}

		if errors.Is(err, storage.ErrSaveOrderAlreadyExists) {
			return response.Status(http.StatusOK)
		}

		return response.Status(http.StatusInternalServerError)
	}

	// TODO: Figure out how to handle accruals
	//accrual, err := s.Accrual.GetAccrual(orderID)
	//if err != nil {
	//	return response.Status(http.StatusInternalServerError)
	//}

	return response.Status(http.StatusAccepted)
}
