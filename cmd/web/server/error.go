package server

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

func (app *Application) reportServerError(r *http.Request, err error) {
	var (
		message = err.Error()
		method  = r.Method
		url     = r.URL.String()
		trace   = string(debug.Stack())
	)

	requestAttrs := slog.Group("request", "method", method, "url", url)
	app.Logger.Error(message, requestAttrs, "trace", trace)
}

func (app *Application) ServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.reportServerError(r, err)

	message := "The server encountered a problem and could not process your request"
	http.Error(w, message, http.StatusInternalServerError)
}

func (app *Application) NotFound(w http.ResponseWriter, r *http.Request) {
	message := "The requested resource could not be found"
	http.Error(w, message, http.StatusNotFound)
}

func (app *Application) BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	message := "Bad request: " + err.Error()
	http.Error(w, message, http.StatusBadRequest)
}
