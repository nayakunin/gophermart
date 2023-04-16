package server

import (
	"io"
	"net/http"
	"strconv"

	api "github.com/nayakunin/gophermart/internal/generated"
)

func (s Server) PostAPIUserOrders(_ http.ResponseWriter, r *http.Request) *api.Response {
	response := api.Response{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return response.Status(http.StatusInternalServerError)
	}

	orderID, err := strconv.Atoi(string(body))
	if err != nil {
		return response.Status(http.StatusBadRequest)
	}

	err = s.Storage.SaveOrder(int64(orderID), r.Context().Value("login").(string))
	if err != nil {
		return nil
	}

	return response.Status(http.StatusOK)
}
