package main

import (
	"net/http"
	"os"
	"path/filepath"
)

func doesFileExist(filePath string) bool {
  info, err := os.Stat(filePath)
  if err != nil || info.IsDir() {
    return false
  }
  return true
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type","text/html; charset=utf-8")

  if r.URL.Path == "/" {
    r.URL.Path = "index"
  }

  if !doesFileExist(filepath.Join("html", r.URL.Path + ".html")) {
    http.Error(w, "Page not found", http.StatusNotFound)
    return
  }
  
  http.ServeFile(w, r, filepath.Join("html", r.URL.Path + ".html"))
}

func examHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type","text/html; charset=utf-8")

}
