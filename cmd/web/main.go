package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

type application struct {
	infoLog, errorLog *log.Logger
}

func main() {
	addr := flag.Int("addr", 4000, "HTTP network address")
	static := flag.String("static-dir", "./ui/static", "Static files directory")
	flag.Parse()

	app := &application{
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)
	mux.HandleFunc("/", app.home)

	fileServer := http.FileServer(http.Dir(*static))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	app.infoLog.Println("Starting server on port", *addr)

	srv := &http.Server{
		Addr:     fmt.Sprint(":", *addr),
		ErrorLog: app.errorLog,
		Handler:  mux,
	}

	err := srv.ListenAndServe()
	app.errorLog.Fatal("ERROR", err)
}
