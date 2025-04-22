package task

import (
	"context"
	"github.com/google/uuid"
	"sync"
	"time"
)

type Status string

const (
	Completed  Status = "completed"
	Failed     Status = "failed"
	Processing Status = "processing"
	Waiting    Status = "waiting"
)

type Task struct {
	ID       string
	Status   Status
	Output   string
	Error    error
	Start    time.Time
	End      time.Time
	function func(ctx context.Context) (string, error)
}

type Manager struct {
	mutex   sync.Mutex
	tasks   map[string]*Task
	jobs    chan *Task
	workNum int
}

func (manager *Manager) CreateTask(function func(ctx context.Context) (string, error)) string {
	id := uuid.NewString()
	task := &Task{
		ID:       id,
		Status:   Waiting,
		Start:    time.Now(),
		function: function,
	}
	manager.mutex.Lock()
	manager.tasks[id] = task
	manager.mutex.Unlock()

	manager.jobs <- task

	return id
}

func (manager *Manager) Get(id string) (*Task, bool) {
	manager.mutex.Lock()
	task, found := manager.tasks[id]
	manager.mutex.Unlock()
	return task, found
}

func NewManager(workCount int) *Manager {
	manager := &Manager{
		tasks:   make(map[string]*Task),
		jobs:    make(chan *Task, workCount),
		workNum: workCount}
	return manager
}

func (manager *Manager) Init(ctx context.Context) {
	for i := 0; i < manager.workNum; i++ {
		go manager.worker(ctx)
	}
}

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
}
