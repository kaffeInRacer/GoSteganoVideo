package routes

import (
	handler2 "kaffein/cmd/web/handler"
	"kaffein/cmd/web/middleware"
	"log/slog"
	"net/http"
)

func Routes(logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	mw := middleware.NewMiddlewareHandler(logger)

	mux.Handle("/assets/", http.HandlerFunc(handler2.Resources))
	mux.HandleFunc("/d/", handler2.DownloadEncryptFile)
	mux.HandleFunc("/", handler2.IndexEncode)
	mux.HandleFunc("/decrypt", handler2.IndexDecode)

	return mw.LogAccess(
		mw.RecoverPanic(
			middleware.SecurityHeaders(
				middleware.GzipCompression(mux),
			),
		),
	)
}
