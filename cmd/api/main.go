package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"http_io_bound/config"
	"http_io_bound/internal/task"
	"http_io_bound/internal/web"
)

func main() {
	set, err := config.Load()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	workerCount := set.Task.WorkerCount
	manager := task.NewManager(workerCount)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	manager.Init(ctx)

	router := chi.NewRouter()
	router.Use(web.Logging)
	router.Use(web.Recover)

	handler := web.NewHandler(manager)
	handler.Routes(router)

	srv := &http.Server{
		Addr:         ":" + set.API.Port,
		Handler:      router,
		ReadTimeout:  set.API.ReadTimeout,
		WriteTimeout: set.API.WriteTimeout,
		IdleTimeout:  set.API.IdleTimeout,
	}

	go func() {
		log.Printf("Server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutdown signal received, stopping server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server stopped")
}
