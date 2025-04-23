package web

import (
	"context"
	"encoding/json"
	"fmt"
	"http_io_bound/internal/errlog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"http_io_bound/internal/task"
)

type Handler struct {
	Manager *task.Manager
}

type TaskRequest struct {
	Type   string          `json:"type"`
	Params json.RawMessage `json:"params"`
}

func NewHandler(manager *task.Manager) *Handler {
	return &Handler{Manager: manager}
}

func (handler *Handler) Routes(router chi.Router) {
	router.Post("/tasks", handler.CreateTask)
	router.Get("/tasks/{id}", handler.GetStatus)
	router.Get("/tasks/{id}/result", handler.GetResult)
	router.Get("/tasks", handler.ListTasks)
	router.Get("/tasks/clear", handler.Clear)
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

}

func (handler *Handler) CreateTask(writer http.ResponseWriter, request *http.Request) {
	errlog.Post("Got JSON, Decoding...")
	var taskrequest TaskRequest
	if err := json.NewDecoder(request.Body).Decode(&taskrequest); err != nil {
		errlog.HTTPCheck(writer, "invalid JSON payload", http.StatusBadRequest)
		return
	}

	errlog.Post("Creating Task...")
	var fn func(context.Context) (string, error)
	switch taskrequest.Type {
	case "stub":
		fn = task.IoTask
	default:
		errlog.HTTPCheck(writer, "unknown task type", http.StatusBadRequest)
		return
	}

	id := handler.Manager.CreateTask(fn)
	errlog.Post(fmt.Sprintf("Created Task ID: %s", id))

	writer.Header().Set("Location", "/tasks/"+id)
	writer.WriteHeader(http.StatusAccepted)
}

func (handler *Handler) GetStatus(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	tsk, ok := handler.Manager.Get(id)
	if !ok {
		errlog.HTTPCheck(writer, "task not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"id":     tsk.ID,
		"status": tsk.Status,
		"start":  tsk.Start,
		"end":    tsk.End,
	}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(response)
}

func (handler *Handler) GetResult(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	tsk, ok := handler.Manager.Get(id)
	if !ok {
		errlog.HTTPCheck(writer, "task not found", http.StatusNotFound)
		return
	}

	switch tsk.Status {
	case task.Waiting, task.Processing:
		errlog.HTTPCheck(writer, "task not finished", http.StatusConflict)
		return
	case task.Failed:
		errlog.HTTPCheck(writer, tsk.Error.Error(), http.StatusInternalServerError)
		return
	case task.Completed:
		response := map[string]string{"result": tsk.Output}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(response)
		return
	default:
		errlog.HTTPCheck(writer, "unknown status", http.StatusInternalServerError)
		return
	}
}

func (handler *Handler) ListTasks(w http.ResponseWriter, r *http.Request) {
	tasks := handler.Manager.ListTasks()

	type taskInfo struct {
		ID     string    `json:"id"`
		Status string    `json:"status"`
		Start  time.Time `json:"start"`
	}

	resp := make([]taskInfo, len(tasks))
	for i, t := range tasks {
		resp[i] = taskInfo{
			ID:     t.ID,
			Status: string(t.Status),
			Start:  t.Start,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (handler *Handler) Clear(writer http.ResponseWriter, request *http.Request) {
	errlog.Post("Clearing old tasks")
	handler.Manager.RemoveOldTasks(15 * time.Minute)
	errlog.Post("Old tasks cleared")
	errlog.HTTPCheck(writer, "Old tasks cleared", http.StatusOK)
}
