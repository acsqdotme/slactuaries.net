package main

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (grw gzipResponseWriter) Write(data []byte) (int, error) {
	return grw.Writer.Write(data)
}

func gzipHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		gzipWriter := gzip.NewWriter(w)
		defer gzipWriter.Close()
		gzippedResponseWriter := gzipResponseWriter{Writer: gzipWriter, ResponseWriter: w}
		next.ServeHTTP(gzippedResponseWriter, r)
	})
}

func redirectWWW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.Host, "www.") {
			http.Redirect(w, r, "http://"+strings.TrimPrefix(r.Host, "www.")+r.RequestURI, 302) // TODO make scheme dynamic with flag
			return
		}

		next.ServeHTTP(w, r)
	})
}

func redirectHandler(url string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, url, 302)
	}
}
