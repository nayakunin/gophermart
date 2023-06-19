package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/nayakunin/gophermart/internal/config"
	api "github.com/nayakunin/gophermart/internal/generated"
	"github.com/nayakunin/gophermart/internal/logger"
	"github.com/nayakunin/gophermart/internal/storage"
)

func TestPostAPIUserOrders(t *testing.T) {
	userID := int64(1)
	type args struct {
		body               string
		storageResponseErr error
	}

	tests := []struct {
		name string
		args args
		want *api.Response
	}{
		{
			name: "OK",
			args: args{
				body: `{"order": "1", "sum": 1}`,
			},
			want: &api.Response{Code: http.StatusOK},
		},
		{
			name: "Unable to read body",
			want: &api.Response{Code: http.StatusBadRequest},
		},
		{
			name: "Unable to unmarshal body",
			args: args{
				body: `{"order": "1", "sum": 1`,
			},
			want: &api.Response{Code: http.StatusBadRequest},
		},
		{
			name: "No orderID",
			args: args{
				body: `{"order": "", "sum": 1}`,
			},
			want: &api.Response{Code: http.StatusBadRequest},
		},
		{
			name: "No Sum",
			args: args{
				body: `{"order": "1", "sum": 0}`,
			},
			want: &api.Response{Code: http.StatusBadRequest},
		},
		{
			name: "OrderID is not a number",
			args: args{
				body: `{"order": "abc", "sum": 1}`,
			},
			want: &api.Response{Code: http.StatusBadRequest},
		},
		{
			name: "Order not found",
			args: args{
				body:               `{"order": "1", "sum": 1}`,
				storageResponseErr: storage.ErrWithdrawOrderNotFound,
			},
			want: &api.Response{Code: http.StatusUnprocessableEntity},
		},
		{
			name: "Not enough balance",
			args: args{
				body:               `{"order": "1", "sum": 1}`,
				storageResponseErr: storage.ErrWithdrawBalanceNotEnough,
			},
			want: &api.Response{Code: http.StatusPaymentRequired},
		},
		{
			name: "Unknown storage error",
			args: args{
				body:               `{"order": "1", "sum": 1}`,
				storageResponseErr: fmt.Errorf("error"),
			},
			want: &api.Response{Code: http.StatusInternalServerError},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Config{
				AuthKey: "auth_key",
			}

			storageParams := StorageParams{}

			if tt.args.storageResponseErr != nil {
				storageParams.WithdrawError = tt.args.storageResponseErr
			}

			mockStorage := NewMockStorage(storageParams)
			server := NewMockServer(nil, mockStorage, cfg, nil)
			logger.Mock()

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/user/balance/withdraw", strings.NewReader(tt.args.body))
			r = r.WithContext(context.WithValue(r.Context(), cfg.AuthKey, userID))

			if got := server.PostAPIUserBalanceWithdraw(w, r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostAPIUserBalanceWithdraw() = %v, want %v", got, tt.want)
			}
		})
	}
}
