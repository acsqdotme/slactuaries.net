package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":4005", "Server Port Address")
	flag.Parse()

	mux := http.NewServeMux()

	mux.HandleFunc("/", pageHandler)
	mux.HandleFunc("/p", topicHandler)
	mux.HandleFunc("/p/", topicHandler)
	mux.HandleFunc("/fm", topicHandler)
	mux.HandleFunc("/fm/", topicHandler)
	// todo make file server handler func
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, gzipHandler(redirectWWW(mux)))
	log.Fatal(err)
}
