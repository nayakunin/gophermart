package server

import (
	"net/http"
	"strconv"

	api "github.com/nayakunin/gophermart/internal/generated"
)

func (s Server) GetAPIUserOrders(_ http.ResponseWriter, r *http.Request) *api.Response {
	userID := r.Context().Value("login").(string)

	orders, err := s.Storage.GetOrders(userID)
	if err != nil {
		return nil
	}

	var apiOrders []api.GetOrdersOrder
	for _, order := range orders {
		apiOrders = append(apiOrders, api.GetOrdersOrder{
			Accrual:   nil,
			Number:    strconv.FormatInt(order.ID, 10),
			Status:    api.OrderStatus{},
			UpdatedAt: order.UploadedAt,
		})
	}

	return api.GetAPIUserOrdersJSON200Response(api.GetOrdersReply{
		Orders: apiOrders,
	})
}
