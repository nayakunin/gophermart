package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/nayakunin/gophermart/internal/config"
	api "github.com/nayakunin/gophermart/internal/generated"
	"github.com/nayakunin/gophermart/internal/logger"
	"github.com/nayakunin/gophermart/internal/storage"
)

func TestGetAPIUserWithdrawals(t *testing.T) {
	userID := int64(1)
	type args struct {
		transactions []storage.Transaction
		err          error
	}

	processedAt := time.Time{}

	tests := []struct {
		name string
		args args
		want *api.Response
	}{
		{
			name: "OK",
			args: args{
				transactions: []storage.Transaction{{
					ID:          1,
					UserID:      userID,
					OrderID:     1,
					Amount:      1,
					ProcessedAt: processedAt,
				}},
			},
			want: api.GetAPIUserWithdrawalsJSON200Response([]api.GetUserWithdrawalsReplyItem{{
				Sum:         1,
				Order:       "1",
				ProcessedAt: processedAt,
			}}),
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

			if tt.args.transactions != nil {
				storageParams.GetWithdrawalsResponse = tt.args.transactions
			}

			if tt.args.err != nil {
				storageParams.GetWithdrawalsError = tt.args.err
			}

			mockStorage := NewMockStorage(storageParams)
			server := NewMockServer(nil, mockStorage, cfg, nil)
			logger.Mock()

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/user/withdrawals", nil)
			r = r.WithContext(context.WithValue(r.Context(), cfg.AuthKey, userID))

			if got := server.GetAPIUserWithdrawals(w, r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAPIUserWithdrawals() = %v, want %v", got, tt.want)
			}
		})
	}
}
