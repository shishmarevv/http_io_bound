package task

import (
	"context"
	"fmt"
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
