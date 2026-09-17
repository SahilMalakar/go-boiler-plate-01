package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {

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

	// Creates the HTTP server and configures how it accepts and handles requests.
	srv := http.Server{
		Addr:         ":8090",
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	// Starts the server and keeps the application running while it accepts requests.
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
