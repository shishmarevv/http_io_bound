package task

import (
	"context"
	"fmt"
	"http_io_bound/internal/errlog"
	"io"
	"net/http"
	"time"

	"http_io_bound/config"
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

func IoTask(ctx context.Context) (string, error) {
	set, err := config.Load()
	errlog.Check("Can't load config", err, true)
	request, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:"+set.IOserver.Port+"/process", nil)
	if err != nil {
		return "", err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %s", response.Status)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
