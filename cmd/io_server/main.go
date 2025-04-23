package main

import (
	"encoding/json"
	"errors"
	"http_io_bound/config"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type response struct {
	Message   string        `json:"message"`
	Duration  time.Duration `json:"duration_ms"`
	Timestamp time.Time     `json:"timestamp"`
}

func processHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	delay := time.Duration(3*60+rand.Intn(2*60)) * time.Second
	log.Printf("New request, delay: %v\n", delay)

	select {
	case <-time.After(delay):
		response := response{
			Message:   "Task completed successfully",
			Duration:  delay / time.Millisecond,
			Timestamp: time.Now(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	case <-ctx.Done():
		log.Println("Request cancelled by client")
		http.Error(w, "request cancelled", http.StatusRequestTimeout)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/process", processHandler)

	set, err := config.Load()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	server := &http.Server{
		Addr:         ":" + set.IOserver.Port,
		Handler:      mux,
		ReadTimeout:  set.IOserver.ReadTimeout,
		WriteTimeout: set.IOserver.WriteTimeout,
		IdleTimeout:  set.IOserver.IdleTimeout,
	}

	log.Printf("I/O-bound server listening on %v\n", set.IOserver.Port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("I/o-bound Server failed: %v", err)
	}
}
