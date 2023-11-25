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

	"slactuaries.net/pkgs/direr"
)

const (
	htmlDir     = "./html/"
	tmplFileExt = ".tmpl.html"
)

func doesFileExist(pathToFile string) bool {
	info, err := os.Stat(filepath.Clean(pathToFile))
	if err != nil || info.IsDir() {
		return false
	}
	return true
}

func doesDirExist(pathToFile string) bool {
	info, err := os.Stat(filepath.Clean(pathToFile))
	if err != nil || !info.IsDir() {
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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	path := strings.Split(r.URL.Path, "/")
	page := path[1]
	if r.URL.Path == "/" {
		page = "index"
	}

	if len(path) == 3 && path[2] == "" {
		http.Redirect(w, r, "/"+page, 302)
		return
	} else if len(path) > 3 {
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}

	if !doesFileExist(filepath.Join(htmlDir, "pages", page+tmplFileExt)) {
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}

	tmpl, err := bindTMPL(
		filepath.Join(htmlDir, "base"+tmplFileExt),
		filepath.Join(htmlDir, "pages", page+tmplFileExt),
	)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	serveTMPL(w, r, tmpl, nil)
}

func topicHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	path := strings.Split(r.URL.Path, "/")
	topic := path[1]

	if !doesDirExist(filepath.Join(htmlDir, "topics", topic)) {
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}

	if len(path) == 3 && path[2] == "" {
		http.Redirect(w, r, "/"+topic, 302)
		return
	}

	data := make(map[string]any)
	t, err := direr.GenerateTree(filepath.Join(htmlDir, "topics", topic, "lessons"), tmplFileExt, filepath.Join("/", topic)) //TODO change to tmpl html
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	data["Tree"] = *t

	if len(path) < 3 {
		tmpl, err := bindTMPL(
			filepath.Join(htmlDir, "base"+tmplFileExt),
			filepath.Join(htmlDir, "topics", topic, "index"+tmplFileExt),
			filepath.Join(htmlDir, "partials", "ulTree"+tmplFileExt),
		)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		serveTMPL(w, r, tmpl, data)
		return
	}

	page := strings.TrimPrefix(r.URL.Path, "/"+topic+"/") // TODO make this less ugly but still removing leading slash

	if !doesFileExist(filepath.Join(htmlDir, "topics", topic, "lessons", page+tmplFileExt)) { // TODO change how this checks for files
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}

	// change tree to subtree, if exists
	t = direr.GetSubTree(t, filepath.Clean(r.URL.Path))
	if t == nil {
		log.Println(errors.New("subtree meta data not found"))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	data["Tree"] = *t

	tmpl, err := bindTMPL(
		filepath.Join(htmlDir, "base"+tmplFileExt),
		filepath.Join(htmlDir, "topics", topic, "lessons", page+tmplFileExt),
	)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	serveTMPL(w, r, tmpl, data)
	return
}
