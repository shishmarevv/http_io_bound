package web

import (
	"context"
	"encoding/json"
	"http_io_bound/internal/errlog"
	"net/http"

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
}

func (handler *Handler) CreateTask(writer http.ResponseWriter, request *http.Request) {
	var taskrequest TaskRequest
	if err := json.NewDecoder(request.Body).Decode(&taskrequest); err != nil {
		errlog.HTTPCheck(writer, "invalid JSON payload", http.StatusBadRequest)
		return
	}

	var fn func(context.Context) (string, error)
	switch taskrequest.Type {
	case "stub":
		fn = task.IoTask
	default:
		errlog.HTTPCheck(writer, "unknown task type", http.StatusBadRequest)
		return
	}

	id := handler.Manager.CreateTask(fn)

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

func (handler *Handler) GetResult(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tsk, ok := handler.Manager.Get(id)
	if !ok {
		errlog.HTTPCheck(w, "task not found", http.StatusNotFound)
		return
	}

	switch tsk.Status {
	case task.Waiting, task.Processing:
		errlog.HTTPCheck(w, "task not finished", http.StatusConflict)
		return
	case task.Failed:
		errlog.HTTPCheck(w, tsk.Error.Error(), http.StatusInternalServerError)
		return
	case task.Completed:
		resp := map[string]string{"result": tsk.Output}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	default:
		errlog.HTTPCheck(w, "unknown status", http.StatusInternalServerError)
		return
	}
}
