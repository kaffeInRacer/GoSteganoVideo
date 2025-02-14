package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type GzipResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (w GzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GzipCompression(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Cek apakah client menerima gzip encoding
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()

		gzipWriter := GzipResponseWriter{ResponseWriter: w, Writer: gz}
		next.ServeHTTP(gzipWriter, r)
	})
}
