package server

import (
	"net/http"
	"time"

	api "github.com/nayakunin/gophermart/internal/generated"
)

func (Server) GetAPIUserWithdrawals(w http.ResponseWriter, r *http.Request) *api.Response {
	return api.GetAPIUserWithdrawalsJSON200Response(api.GetUserWithdrawalsReply{{
		Order:       "",
		ProcessedAt: time.Time{},
		Sum:         0,
	}})
}
