package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/fntlnz/netcan/ws"
	"github.com/gorilla/handlers"
	"github.com/sirupsen/logrus"
)

var rootfsVar string

func main() {
	flag.StringVar(&rootfsVar, "rootfs", "", "Specify a rootfs where to get network information")
	flag.Parse()
	loggedRouter := handlers.LoggingHandler(os.Stdout, ws.NewRouter(rootfsVar))

	srv := &http.Server{
		Handler:      loggedRouter,
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logrus.Fatal(srv.ListenAndServe())
}
