package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seka/reci-pin/backend/internal/server/middleware"
	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mh := middleware.CORS(nextHandler)

	t.Run("Headers are set correctly for GET", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		mh.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
		assert.Equal(t, "Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
		assert.Equal(t, "86400", w.Header().Get("Access-Control-Max-Age"))
		assert.Equal(t, "OK", w.Body.String())
	})

	t.Run("OPTIONS Request yields No Content and specific headers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/", nil)
		w := httptest.NewRecorder()

		mh.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
		assert.Equal(t, "Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
		assert.Equal(t, "86400", w.Header().Get("Access-Control-Max-Age"))
		// Body should be empty for 204
		assert.Empty(t, w.Body.String())
	})
}
