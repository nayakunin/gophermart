package server

import (
	"net/http"

	api "github.com/nayakunin/gophermart/internal/generated"
)

func (s Server) GetAPIUserBalance(_ http.ResponseWriter, r *http.Request) *api.Response {
	userID := r.Context().Value("login").(string)

	current, withdrawn, err := s.Storage.GetBalance(userID)
	if err != nil {
		return nil
	}

	return api.GetAPIUserBalanceJSON200Response(api.Balance{
		Current:   current,
		Withdrawn: withdrawn,
	})
}
