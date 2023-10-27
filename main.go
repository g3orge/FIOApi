package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/g3orge/FIOApi/internal/db"
	"github.com/g3orge/FIOApi/internal/transport"
	"github.com/gorilla/mux"
)

func main() {
	db.InitPostgres()
	r := mux.NewRouter()
	r.HandleFunc("/get{id}", transport.GetF).Methods("GET")
	r.HandleFunc("/update{id}", transport.UpdateF).Methods("PUT")
	r.HandleFunc("/delete{id}", transport.DeleteF).Methods("DELETE")
	r.HandleFunc("/add", transport.AddF).Methods("POST")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	stopped := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}
		close(stopped)
	}()

	log.Printf("Starting HTTP server on %s", "cfg.HTTPAddr")

	// start HTTP server
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	}

	<-stopped

	log.Printf("Have a nice day!")
}
