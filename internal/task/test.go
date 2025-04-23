package task

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestLifecycle(test *testing.T) {
	manager := NewManager(1)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	manager.Init(ctx)

	id := manager.CreateTask(func(ctx context.Context) (string, error) {
		return "ok", nil
	})

	task, found := manager.Get(id)
	if !found {
		test.Fatalf("Task %s not found after CreateTask", id)
	}
	if task.Status != Waiting && task.Status != Processing {
		test.Errorf("Expected pending or running, got %s", task.Status)
	}

	time.Sleep(20 * time.Millisecond)

	task, _ = manager.Get(id)
	if task.Status != Completed {
		test.Errorf("Expected completed, got %s", task.Status)
	}
	if task.Output != "ok" {
		test.Errorf("Expected result 'ok', got '%s'", task.Output)
	}
	if task.Error != nil {
		test.Errorf("Expected no error, got %v", task.Error)
	}
}

func TestError(test *testing.T) {
	manager := NewManager(1)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	manager.Init(ctx)

	id := manager.CreateTask(func(ctx context.Context) (string, error) {
		return "", errors.New("fail")
	})

	time.Sleep(20 * time.Millisecond)
	task, _ := manager.Get(id)
	if task.Status != Failed {
		test.Errorf("Expected failed, got %s", task.Status)
	}
	if task.Error == nil || task.Error.Error() != "fail" {
		test.Errorf("Expected error 'fail', got %v", task.Error)
	}
}
