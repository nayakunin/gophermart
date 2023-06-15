package server

import (
	"net/http"
	"strconv"

	api "github.com/nayakunin/gophermart/internal/generated"
	"github.com/nayakunin/gophermart/internal/logger"
	"github.com/nayakunin/gophermart/internal/middlewares"
)

func (s Server) GetAPIUserWithdrawals(_ http.ResponseWriter, r *http.Request) *api.Response {
	response := api.Response{}
	userID := r.Context().Value(middlewares.AuthKey).(int64)

	withdrawals, err := s.Storage.GetWithdrawals(userID)
	if err != nil {
		logger.Errorf("failed to get withdrawals: %v", err)
		return response.Status(http.StatusInternalServerError)
	}

	var apiWithdrawals []api.GetUserWithdrawalsReplyItem
	for _, withdrawal := range withdrawals {
		apiWithdrawals = append(apiWithdrawals, api.GetUserWithdrawalsReplyItem{
			Order:       strconv.FormatInt(withdrawal.OrderID, 10),
			ProcessedAt: withdrawal.ProcessedAt,
			Sum:         withdrawal.Amount,
		})
	}

	return api.GetAPIUserWithdrawalsJSON200Response(apiWithdrawals)
}
