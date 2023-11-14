package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
  addr := flag.String("addr", ":4004", "Server Port Address")
  flag.Parse()

  mux := http.NewServeMux()
  mux.HandleFunc("/", pageHandler)
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

  log.Printf("Starting server on %s", *addr)
  err := http.ListenAndServe(*addr, mux)
  log.Fatal(err)
}
