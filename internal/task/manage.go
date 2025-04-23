package task

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"http_io_bound/internal/errlog"
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
	mutex   sync.RWMutex
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
	manager.mutex.RLock()

	errlog.Post(fmt.Sprintf("Looking for task %s", id))
	task, found := manager.tasks[id]
	manager.mutex.RUnlock()
	return task, found
}

func NewManager(workCount int) *Manager {
	manager := &Manager{
		tasks:   make(map[string]*Task),
		jobs:    make(chan *Task, workCount*3),
		workNum: workCount}
	return manager
}

func (manager *Manager) Init(ctx context.Context) {
	for i := 0; i < manager.workNum; i++ {
		go manager.worker(ctx)
	}
}

func (manager *Manager) ListTasks() []*Task {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	list := make([]*Task, 0, len(manager.tasks))
	for _, t := range manager.tasks {
		list = append(list, t)
	}
	return list
}

func (manager *Manager) RemoveOldTasks(duration time.Duration) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	for id, task := range manager.tasks {
		if time.Since(task.End) > duration {
			delete(manager.tasks, id)
		}
	}
}
