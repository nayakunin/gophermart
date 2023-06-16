package server

import (
	"net/http"

	api "github.com/nayakunin/gophermart/internal/generated"
	"github.com/nayakunin/gophermart/internal/logger"
)

func (s Server) GetAPIUserBalance(_ http.ResponseWriter, r *http.Request) *api.Response {
	response := api.Response{}

	userID := r.Context().Value(s.Cfg.AuthKey).(int64)

	current, withdrawn, err := s.Storage.GetBalance(userID)
	if err != nil {
		logger.Errorf("failed to get balance: %v", err)
		return response.Status(http.StatusInternalServerError)
	}

	return api.GetAPIUserBalanceJSON200Response(api.Balance{
		Current:   current,
		Withdrawn: withdrawn,
	})
}
