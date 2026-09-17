package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sahil_malakar/production_grade_golang_setup/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println("go app server is running!!..")

	// Creates a router that receives incoming requests
	// and forwards them to the handler registered for the matching route.
	mux := http.NewServeMux()

	// Registers the GET /health route and executes this function
	// whenever a request is made to that endpoint.
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	log.Printf("Server is listening on http://localhost:%v", cfg.Port)

	// Creates the HTTP server and configures how it accepts and handles requests.
	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// Starts the server and keeps the application running while it accepts requests.
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
