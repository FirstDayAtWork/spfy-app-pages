package main

import (
	"fmt"
	"log"
	"mustracker/basic_server"
	"net/http"
)

const port int = 2228

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/register":
		switch r.Method {
		case "GET":
			basic_server.RenderRegister(w, r)
		case "POST":
			// TODO implement parsing body and writing to DB
			fmt.Fprintf(w, "%s is called for %s", r.Method, r.URL.Path)
		default:
			fmt.Fprintf(w, "ERROR! %s is not supported for %s", r.Method, r.URL.Path)
		}
	}
}

func main() {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fs))
	mux.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
