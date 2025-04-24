package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"http_io_bound/config"
	"http_io_bound/internal/errlog"
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

func processHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	delay := time.Duration(3*60+rand.Intn(2*60)) * time.Second
	errlog.Post(fmt.Sprintf("[STUB] New request, delay: %v", delay))

	select {
	case <-time.After(delay):
		response := response{
			Message:   "Task completed successfully",
			Duration:  delay / time.Millisecond,
			Timestamp: time.Now(),
		}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(response)
	case <-ctx.Done():
		errlog.Post("[STUB] Request cancelled by client")
		errlog.HTTPCheck(writer, "[STUB] request cancelled", http.StatusRequestTimeout)
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

	errlog.Post(fmt.Sprintf("[STUB] I/O-bound server listening on %v\n", set.IOserver.Port))
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		errlog.Check("[STUB] I/o-bound Server failed", err, true)
	}
}
