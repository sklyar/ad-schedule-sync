package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckSecretKeyMiddleware(t *testing.T) {
	t.Parallel()

	const secretKey = "secret key"
	mw := checkSecretKeyMiddleware(secretKey)

	tests := []struct {
		name       string
		headers    map[string]string
		wantStatus int
	}{
		{
			name: "valid key",
			headers: map[string]string{
				secretKeyHeader: secretKey,
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "invalid key",
			headers: map[string]string{
				secretKeyHeader: "invalid",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "no key",
			headers:    map[string]string{},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequest(http.MethodGet, "/", nil)
			require.NoError(t, err)

			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}

			handler := func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}

			w := httptest.NewRecorder()
			mw(http.HandlerFunc(handler)).ServeHTTP(w, req)

			require.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
