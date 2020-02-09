package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request")
	w.Write([]byte("Hello, world"))
}

func waitForShutdown(server *http.Server) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGINT,
		syscall.SIGTERM)

	// Block until we receive our signal.
	_ = <-sig
	log.Printf("API server shutting down")
	server.Shutdown(context.Background())
	log.Printf("server shutdown complete")
}

func main() {
	//new router instance.
	r := mux.NewRouter()

	// register routes mapping URL paths to handlers
	r.HandleFunc("/", handler).Methods(http.MethodGet)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:9000",
		WriteTimeout: 15 * time.Second, // it is good practice for enforce timeout for server
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		log.Printf("Starting http server at - %s", srv.Addr)
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatalf("Error when starting the server = %s", err.Error())
		}
	}()

	// Do Graceful Shutdown
	waitForShutdown(srv)
}
