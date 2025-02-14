package handler

import (
	"kaffein/assets/static"
	"net/http"
	"strings"
)

func Resources(w http.ResponseWriter, r *http.Request) {
	staticFileServer := http.FileServer(http.FS(static.StaticFiles))
	path := strings.TrimPrefix(r.URL.Path, "/assets/")
	if strings.HasSuffix(path, "/") || path == "" {
		http.NotFound(w, r)
		return
	}
	http.StripPrefix("/assets/", staticFileServer).ServeHTTP(w, r)
}
