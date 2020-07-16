package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	addr := flag.Int("addr", 4000, "HTTP network address")
	static := flag.String("static-dir", "./ui/static", "Static files directory")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)
	mux.HandleFunc("/", home)

	fileServer := http.FileServer(http.Dir(*static))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Println("Starting server on port", *addr)
	err := http.ListenAndServe(fmt.Sprint(":", *addr), mux)
	log.Fatal(err)
}
