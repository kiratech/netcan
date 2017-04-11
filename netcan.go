package main

import (
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/fntlnz/netcan/ws"
	"github.com/gorilla/handlers"
)

func main() {
	loggedRouter := handlers.LoggingHandler(os.Stdout, ws.NewRouter())

	srv := &http.Server{
		Handler:      loggedRouter,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logrus.Fatal(srv.ListenAndServe())
}
