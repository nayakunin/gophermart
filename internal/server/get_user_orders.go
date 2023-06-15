package server

import (
	"net/http"
	"strconv"

	api "github.com/nayakunin/gophermart/internal/generated"
	"github.com/nayakunin/gophermart/internal/logger"
	"github.com/nayakunin/gophermart/internal/middlewares"
)

func (s Server) GetAPIUserOrders(_ http.ResponseWriter, r *http.Request) *api.Response {
	response := api.Response{}
	userID := r.Context().Value(middlewares.AuthKey).(int64)

	orders, err := s.Storage.GetOrders(userID)
	if err != nil {
		logger.Errorf("failed to get orders: %v", err)
		return response.Status(http.StatusInternalServerError)
	}

	apiOrders := make([]api.GetOrdersOrder, 0, len(orders))
	for _, order := range orders {
		apiOrders = append(apiOrders, api.GetOrdersOrder{
			Accrual:    order.Accrual,
			Number:     strconv.FormatInt(order.ID, 10),
			Status:     order.Status,
			UploadedAt: order.UploadedAt,
		})
	}

	return api.GetAPIUserOrdersJSON200Response(apiOrders)
}
