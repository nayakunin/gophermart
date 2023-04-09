package server

import (
	"net/http"
	"time"

	api "github.com/nayakunin/gophermart/internal/generated"
)

func (Server) GetAPIUserOrders(w http.ResponseWriter, r *http.Request) *api.Response {
	return api.GetAPIUserOrdersJSON200Response(api.GetOrdersReply{
		Orders: []api.GetOrdersOrder{{
			Accrual:   nil,
			Number:    "",
			Status:    api.OrderStatus{},
			UpdatedAt: time.Time{},
		}},
	})
}
