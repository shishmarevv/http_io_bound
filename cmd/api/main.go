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

	"http_io_bound/internal/task"
	"http_io_bound/internal/web"
)

func main() {
	workerCount := 5
	tm := task.NewManager(workerCount)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	tm.Init(ctx)

	router := chi.NewRouter()
	router.Use(web.Logging)
	router.Use(web.Recover)

	handler := web.NewHandler(tm)
	handler.Routes(router)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Minute,
		IdleTimeout:  1 * time.Minute,
	}

	go func() {
		log.Printf("🚀 Server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("🛑 Shutdown signal received, stopping server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("✅ Server gracefully stopped")
}
