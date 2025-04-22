package task

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

func (manager *Manager) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case task := <-manager.jobs:
			manager.run(ctx, task)
		}
	}
}

func (manager *Manager) run(ctx context.Context, task *Task) {
	defer func() {
		if r := recover(); r != nil {
			manager.mutex.Lock()
			defer manager.mutex.Unlock()
			task.Status = Failed
			task.Error = fmt.Errorf("panic: %v", r)
			task.End = time.Now()
		}
	}()

	manager.mutex.Lock()
	task.Status = Processing
	manager.mutex.Unlock()

	result, err := task.function(ctx)

	manager.mutex.Lock()
	if err != nil {
		task.Status = Failed
		task.Error = err
	} else {
		task.Status = Completed
		task.Output = result
	}
	task.End = time.Now()
	manager.mutex.Unlock()
}

func ioTask(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:9090/process", nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
