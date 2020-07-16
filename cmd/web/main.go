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

	srv := &http.Server{
		Addr:     fmt.Sprint(":", *addr),
		ErrorLog: app.errorLog,
		Handler:  app.routes(static),
	}

	app.infoLog.Println("Starting server on port", *addr)
	err := srv.ListenAndServe()
	app.errorLog.Fatal("ERROR", err)
}
