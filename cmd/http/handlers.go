package main

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func doesFileExist(pathToFile string) bool {
	info, err := os.Stat(filepath.Clean(pathToFile))
	if err != nil || info.IsDir() {
		return false
	}
	return true
}

func bindTMPL(files ...string) (*template.Template, error) {
	for _, checkFile := range files {
		if !doesFileExist(checkFile) {
			return nil, errors.New("Template file missing " + checkFile)
		}
	}

	tmpl, err := template.New("notSureWhatThisDoes").Funcs(nil).ParseFiles(files...)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func serveTMPL(w http.ResponseWriter, r *http.Request, tmpl *template.Template, data map[string]interface{}) {
	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, "base", data)
	if err != nil {
		log.Println(err.Error())
		// fancyErrorHandler(w, r, http.StatusInternalServerError)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError) // TODO add proper error pages
		return
	}
	buf.WriteTo(w)
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
