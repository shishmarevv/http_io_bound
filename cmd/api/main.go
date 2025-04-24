package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"http_io_bound/cmd/api/commands"
	"http_io_bound/internal/errlog"
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
	errlog.Check("Can't read config: ", err, true)

	err = commands.Execute()
	errlog.Check("CLI error: ", err, true)

	manager := task.NewManager(set.Task.WorkerCount)

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
		errlog.Post(fmt.Sprintf("Server listening on %s", srv.Addr))
		err := srv.ListenAndServe()
		errlog.Check("ListenAndServe", err, true)
	}()

	<-ctx.Done()
	errlog.Post("Shutdown signal received, stopping server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = srv.Shutdown(shutdownCtx)
	errlog.Check("Server shutdown failed", err, false)
	errlog.Post("Server stopped")
}
