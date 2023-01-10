package tasks

import (
	"errors"
	"time"
)

type Task struct {
	Name           string
	Description    string
	Complete       bool
	CompletionDate time.Time
}

type TaskManager struct {
	repo TaskRepository
}

type TaskRepository interface {
	Store(t Task) (Task, error)
	LoadAll() []Task
	Update(task Task) Task
	Load(name string) (*Task, error)
	Delete(task Task)
}

func NewTaskManager(t TaskRepository) *TaskManager {
	return &TaskManager{
		repo: t,
	}
}

func (tm *TaskManager) CreateTask(name string, description string) (*Task, error) {
	t, err := tm.repo.Store(Task{
		Name:        name,
		Description: description,
	})
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (tm *TaskManager) GetTask(name string) (*Task, error) {
	return tm.repo.Load(name)
}

func (tm *TaskManager) ListTasks() []Task {
	return tm.repo.LoadAll()
}

func (tm *TaskManager) Complete(task Task) *Task {
	task.CompletionDate = time.Now()
	task.Complete = true
	t := tm.repo.Update(task)
	return &t
}

func (tm *TaskManager) Remove(name string) {
	tm.repo.Delete(Task{Name: name})
}

func (tm *TaskManager) ListCompletedTasks(hrs int) ([]Task, error) {
	if isInvalid(hrs) {
		return nil, errors.New("invalid period selected, must be 24, 12, or 6")
	}
	tasks := tm.repo.LoadAll()
	var ret []Task
	currentTime := time.Now()
	for _, v := range tasks {
		diff := currentTime.Sub(v.CompletionDate)
		if diff < time.Duration(hrs)*time.Hour {
			ret = append(ret, v)
		}
	}

	return ret, nil
}

func isInvalid(hrs int) bool {
	if hrs == 24 || hrs == 12 || hrs == 6 {
		return false
	}
	return true
}
