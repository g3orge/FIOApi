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
	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load("conf.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	loadEnv()
	db_user := os.Getenv("DB_USER")
	db_pass := os.Getenv("DB_PASS")
	db_port := os.Getenv("DB_PORT")
	db_name := os.Getenv("DB_NAME")
	db_ssl := os.Getenv("DB_SSL_MODE")
	serv_port := os.Getenv("SERV_PORT")
	db.InitPostgres(db_user, db_pass, db_port, db_name, db_ssl)

	r := mux.NewRouter()
	r.HandleFunc("/getAll", transport.GetAll).Methods("GET")
	r.HandleFunc("/get/filter", transport.GetF).Methods("GET")
	r.HandleFunc("/update{id}", transport.UpdateF).Methods("PUT")
	r.HandleFunc("/delete{id}", transport.DeleteF).Methods("DELETE")
	r.HandleFunc("/add", transport.AddF).Methods("POST")

	srv := &http.Server{
		Addr:    serv_port,
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

	log.Printf("Starting HTTP server on %s", serv_port)

	// start HTTP server
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	}

	<-stopped

	log.Printf("Have a nice day!")
}
