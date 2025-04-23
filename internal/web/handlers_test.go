package web

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
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

func TestCreateAndGetStatus(test *testing.T) {
	os.Setenv("LOG_DIR", test.TempDir())

	router, _ := setupServer()

	jsonBody := `{"type":"stub","params":{}}`
	request := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(jsonBody))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusAccepted {
		test.Fatalf("Expected 202, got %d", recorder.Code)
	}
	location := recorder.Header().Get("Location")
	if !strings.HasPrefix(location, "/tasks/") {
		test.Fatalf("Invalid Location header: %q", location)
	}

	request2 := httptest.NewRequest(http.MethodGet, location, nil)
	recorder2 := httptest.NewRecorder()
	router.ServeHTTP(recorder2, request2)

	if recorder2.Code != http.StatusOK {
		test.Errorf("Expected 200 on status, got %d", recorder2.Code)
	}
	var st struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(recorder2.Body).Decode(&st); err != nil {
		test.Fatal(err)
	}
	if st.Status != string(task.Waiting) && st.Status != string(task.Processing) {
		test.Errorf("Unexpected status %q", st.Status)
	}
}

func TestGetResult(test *testing.T) {
	os.Setenv("LOG_DIR", test.TempDir())

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
