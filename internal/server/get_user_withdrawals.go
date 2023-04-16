package server

import (
	"net/http"
	"strconv"

	api "github.com/nayakunin/gophermart/internal/generated"
)

func (s Server) GetAPIUserWithdrawals(_ http.ResponseWriter, r *http.Request) *api.Response {
	userID := r.Context().Value("login").(string)

	withdrawals, err := s.Storage.GetWithdrawals(userID)
	if err != nil {
		return nil
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
