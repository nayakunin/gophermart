package server

import (
	"net/http"
	"strconv"

	api "github.com/nayakunin/gophermart/internal/generated"
)

func (s Server) GetAPIUserOrders(_ http.ResponseWriter, r *http.Request) *api.Response {
	response := api.Response{}
	userID := r.Context().Value("userID").(int64)

	orders, err := s.Storage.GetOrders(userID)
	if err != nil {
		return response.Status(http.StatusInternalServerError)
	}

	apiOrders := make([]api.GetOrdersOrder, 0, len(orders))
	for _, order := range orders {
		accrual := float32(0.0)
		apiOrders = append(apiOrders, api.GetOrdersOrder{
			Accrual:    &accrual,
			Number:     strconv.FormatInt(order.ID, 10),
			Status:     order.Status,
			UploadedAt: order.UploadedAt,
		})
	}

	return api.GetAPIUserOrdersJSON200Response(apiOrders)
}
