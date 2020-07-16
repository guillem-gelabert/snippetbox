package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	addr := flag.Int("addr", 4000, "HTTP network address")
	static := flag.String("static-dir", "./ui/static", "Static files directory")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)
	mux.HandleFunc("/", home)

	fileServer := http.FileServer(http.Dir(*static))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	infoLog.Println("Starting server on port", *addr)

	srv := &http.Server{
		Addr:     fmt.Sprint(":", *addr),
		ErrorLog: errorLog,
		Handler:  mux,
	}

	err := srv.ListenAndServe()
	errorLog.Fatal("ERROR", err)
}
