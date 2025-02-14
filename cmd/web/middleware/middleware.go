// In cmd/web/middleware/middleware.go
package middleware

import (
	"fmt"
	"github.com/tomasen/realip"
	"kaffein/cmd/web/server/response"
	"log/slog"
	"net/http"
)

type MiddlewareHandler struct {
	Logger *slog.Logger
}

func NewMiddlewareHandler(logger *slog.Logger) *MiddlewareHandler {
	return &MiddlewareHandler{Logger: logger}
}

func (m *MiddlewareHandler) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				m.Logger.Error("panic recovered", "error", fmt.Errorf("%s", err))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

func (m *MiddlewareHandler) LogAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Use server.NewMetricsResponseWriter here
		mw := response.NewMetricsResponseWriter(w)
		next.ServeHTTP(mw, r)

		var (
			ip     = realip.FromRequest(r)
			method = r.Method
			url    = r.URL.String()
			proto  = r.Proto
		)

		userAttrs := slog.Group("user", "ip", ip)
		requestAttrs := slog.Group("request", "method", method, "url", url, "proto", proto)
		responseAttrs := slog.Group("response", "status", mw.StatusCode, "size", mw.BytesCount)

		m.Logger.Info("access", userAttrs, requestAttrs, responseAttrs)
	})
}
