package main

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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

func serveTMPL(w http.ResponseWriter, r *http.Request, tmpl *template.Template, data map[string]any) {
	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, "base", data)
	if err != nil {
		log.Println(err.Error())
		fancyErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
}

func fetchBaseData(r *http.Request) (data map[string]any) {
	data = make(map[string]any)
	data["Path"] = r.URL.Path
	data["Host"] = r.Host
	return data
}

func fancyErrorHandler(w http.ResponseWriter, r *http.Request, httpCode int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(httpCode)

	tmpl, err := bindTMPL(
		filepath.Join(htmlDir, "base"+tmplFileExt),
		filepath.Join(htmlDir, "partials", "error_meta"+tmplFileExt),
		filepath.Join(htmlDir, "errors", strconv.Itoa(httpCode)+tmplFileExt),
	)
	if err != nil {
		log.Println(err.Error())
		fancyErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", nil) // TODO write to bytes buffer first like in serveTMPL
	if err != nil {
		log.Println(err.Error())
		fancyErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
	return
}

func pageHandler(w http.ResponseWriter, r *http.Request) {

	path := strings.Split(r.URL.Path, "/")
	page := path[1]
	if r.URL.Path == "/" {
		page = "index"
	}

	if len(path) == 3 && path[2] == "" {
		http.Redirect(w, r, "/"+page, 302)
		return
	} else if len(path) > 3 {
		fancyErrorHandler(w, r, http.StatusNotFound)
		return
	}

	if !doesFileExist(filepath.Join(htmlDir, "pages", page+tmplFileExt)) {
		fancyErrorHandler(w, r, http.StatusNotFound)
		return
	}

	data := fetchBaseData(r)

	tmpl, err := bindTMPL(
		filepath.Join(htmlDir, "base"+tmplFileExt),
		filepath.Join(htmlDir, "pages", page+tmplFileExt),
	)
	if err != nil {
		log.Println(err.Error())
		fancyErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	serveTMPL(w, r, tmpl, data)
	return
}

func topicHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	topic := path[1]

	if !doesDirExist(filepath.Join(htmlDir, "topics", topic)) {
		fancyErrorHandler(w, r, http.StatusNotFound)
		return
	}

	if len(path) == 3 && path[2] == "" {
		http.Redirect(w, r, "/"+topic, 302)
		return
	}

	data := fetchBaseData(r)
	t, err := direr.GenerateTree(filepath.Join(htmlDir, "topics", topic, "lessons"), tmplFileExt, filepath.Join("/", topic)) //TODO change to tmpl html
	if err != nil {
		log.Println(err.Error())
		fancyErrorHandler(w, r, http.StatusInternalServerError)
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
			fancyErrorHandler(w, r, http.StatusInternalServerError)
			return
		}

		serveTMPL(w, r, tmpl, data)
		return
	}

	page := strings.TrimPrefix(r.URL.Path, "/"+topic+"/") // TODO make this less ugly but still removing leading slash
	pageFilePath := filepath.Join(htmlDir, "topics", topic, "lessons", page)

	if strings.HasSuffix(page, ".md") && doesFileExist(pageFilePath) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.ServeFile(w, r, filepath.Join(htmlDir, "topics", topic, "lessons", page))
		return
	}

	if strings.HasSuffix(page, ".pdf") && doesFileExist(pageFilePath) {
		w.Header().Set("Content-Type", "application/pdf")
		http.ServeFile(w, r, filepath.Join(htmlDir, "topics", topic, "lessons", page))
		return
	}

	if !doesFileExist(pageFilePath + tmplFileExt) { // TODO change how this checks for files
		fancyErrorHandler(w, r, http.StatusNotFound)
		return
	}

	// change tree to subtree, if exists
	t = direr.GetSubTree(t, filepath.Clean(r.URL.Path))
	if t == nil {
		log.Println(errors.New("subtree meta data not found"))
		fancyErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
	data["Tree"] = *t

	tmpl, err := bindTMPL(
		filepath.Join(htmlDir, "base"+tmplFileExt),
		filepath.Join(htmlDir, "partials", "article"+tmplFileExt),
		filepath.Join(htmlDir, "topics", topic, "lessons", page+tmplFileExt),
	)
	if err != nil {
		log.Println(err.Error())
		fancyErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	serveTMPL(w, r, tmpl, data)
	return
}
