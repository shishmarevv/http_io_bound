package web

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"http_io_bound/internal/task"
)

func setupServer() (*chi.Mux, *task.Manager) {
	manager := task.NewManager(1)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	mux := chi.NewMux()
	handler := NewHandler(manager)
	handler.Routes(mux)
	manager.Init(ctx)
	cancel()
	return mux, manager
}

func TestCreateAndGetStatus(t *testing.T) {
	router, _ := setupServer()

	request := httptest.NewRequest(http.MethodPost, "/tasks", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusAccepted {
		t.Fatalf("Expected 202, got %d", recorder.Code)
	}
	location := recorder.Header().Get("Location")
	if !strings.HasPrefix(location, "/tasks/") {
		t.Fatalf("Invalid Location header: %q", location)
	}

	request2 := httptest.NewRequest(http.MethodGet, location, nil)
	recorder2 := httptest.NewRecorder()
	router.ServeHTTP(recorder2, request2)

	if recorder2.Code != http.StatusOK {
		t.Errorf("Expected 200 on status, got %d", recorder2.Code)
	}
	var st struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(recorder2.Body).Decode(&st); err != nil {
		t.Fatal(err)
	}
	if st.Status != string(task.Waiting) && st.Status != string(task.Processing) {
		t.Errorf("Unexpected status %q", st.Status)
	}
}

func TestGetResult(test *testing.T) {
	router, manager := setupServer()

	id := manager.CreateTask(func(ctx context.Context) (string, error) {
		return "hello", nil
	})

	time.Sleep(20 * time.Millisecond)

	request := httptest.NewRequest(http.MethodGet, "/tasks/"+id+"/result", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		test.Fatalf("Expected 200, got %d", recorder.Code)
	}
	var response struct {
		Result string `json:"result"`
	}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		test.Fatal(err)
	}
	if response.Result != "hello" {
		test.Errorf("Expected 'hello', got %q", response.Result)
	}
}
