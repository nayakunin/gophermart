package server

import (
	"net/http"

	api "github.com/nayakunin/gophermart/internal/generated"
)

func (Server) GetAPIUserBalance(w http.ResponseWriter, r *http.Request) *api.Response {
	return api.GetAPIUserBalanceJSON200Response(api.Balance{
		Current:   0,
		Withdrawn: 0,
	})
}
