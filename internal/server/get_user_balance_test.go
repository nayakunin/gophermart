package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/nayakunin/gophermart/internal/config"
	api "github.com/nayakunin/gophermart/internal/generated"
	"github.com/nayakunin/gophermart/internal/logger"
	"github.com/nayakunin/gophermart/internal/storage"
)

func TestGetAPIUserBalance(t *testing.T) {
	userID := int64(1)
	type args struct {
		balance *storage.Balance
		err     error
	}

	tests := []struct {
		name string
		args args
		want *api.Response
	}{
		{
			name: "OK",
			args: args{
				balance: &storage.Balance{
					Amount:    1,
					Withdrawn: 2,
				},
			},
			want: api.GetAPIUserBalanceJSON200Response(api.Balance{
				Current:   1,
				Withdrawn: 2,
			}),
		},
		{
			name: "Error",
			args: args{
				err: fmt.Errorf("error"),
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

			if tt.args.balance != nil {
				storageParams.GetBalanceResponse = *tt.args.balance
			}

			if tt.args.err != nil {
				storageParams.GetBalanceError = tt.args.err
			}

			mockStorage := NewMockStorage(storageParams)
			server := NewMockServer(nil, mockStorage, cfg, nil)
			logger.Mock()

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/user/balance", nil)
			r = r.WithContext(context.WithValue(r.Context(), cfg.AuthKey, userID))

			if got := server.GetAPIUserBalance(w, r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAPIUserBalance() = %v, want %v", got, tt.want)
			}
		})
	}
}
