package server

import (
	"net/http"

	api "github.com/nayakunin/gophermart/internal/generated"
)

func (Server) PostAPIUserLogin(w http.ResponseWriter, r *http.Request) *api.Response {
	response := api.Response{}
	return response.Status(http.StatusOK)
}
