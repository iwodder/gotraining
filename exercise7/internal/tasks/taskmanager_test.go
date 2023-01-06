package tasks

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CreateTask(t *testing.T) {
	tm := TaskManager{
		repo: &MemRepository{
			tasks: make([]Task, 0, 10),
		},
	}

	task := tm.CreateTask("build app", "create the task manager app")

	assert.Equal(t, "build app", task.Name)
	assert.Equal(t, "create the task manager app", task.Description)
	assert.False(t, task.Completed)
}

func Test_CreateTaskThenListTasks(t *testing.T) {
	tm := TaskManager{
		repo: &MemRepository{
			tasks: make([]Task, 0, 10),
		},
	}
	task := tm.CreateTask("build app", "create the task manager app")

	taskSlice := tm.ListTasks()

	assert.Equal(t, 1, len(taskSlice), "Expected to find one task")
	assert.Contains(t, taskSlice, *task)
}

func Test_CreateTaskThenLoadTask(t *testing.T) {
	tm := TaskManager{
		repo: &MemRepository{
			tasks: make([]Task, 0, 10),
		},
	}
	tm.CreateTask("build app", "create the task manager app")

	task, _ := tm.GetTask("build app")

	assert.NotNil(t, task)
	assert.Equal(t, "build app", task.Name)
	assert.Equal(t, "create the task manager app", task.Description)
}

func Test_CompleteTask(t *testing.T) {
	tm := TaskManager{
		repo: &MemRepository{
			tasks: make([]Task, 0, 10),
		},
	}
	task := tm.CreateTask("build app", "create the task manager app")

	task = tm.Complete(*task)

	assert.True(t, task.Completed)
}

type MemRepository struct {
	tasks []Task
	id    uint64
}

func (m *MemRepository) Store(t Task) Task {
	newTask := Task{
		Name:        t.Name,
		Description: t.Description,
		Completed:   t.Completed,
		ID:          m.nextId(),
	}
	m.tasks = append(m.tasks, newTask)
	return newTask
}

func (m *MemRepository) LoadAll() []Task {
	ret := make([]Task, len(m.tasks))
	copy(ret, m.tasks)
	return ret
}

func (m *MemRepository) Update(task Task) Task {
	for i, v := range m.tasks {
		if v == task {
			m.tasks[i] = task
		}
	}
	return task
}

func (m *MemRepository) LoadTask(name string) (*Task, error) {
	for _, v := range m.tasks {
		if v.Name == name {
			return &Task{
				ID:          v.ID,
				Name:        v.Name,
				Description: v.Description,
				Completed:   v.Completed,
			}, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Unable to locate a task with name %s", name))
}

func (m *MemRepository) nextId() uint64 {
	m.id++
	return m.id
}
