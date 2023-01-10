package tasks

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var tm TaskManager

func setup(t *testing.T) {
	tm = TaskManager{
		repo: &MemRepository{
			tasks: make([]Task, 0, 10),
		},
	}
}

func assertEqualCompletionDate(t *testing.T, date time.Time) {
	year, month, day := time.Now().Date()
	assert.Equalf(t, day, date.Day(), "Expected matching day, wanted %d, got %d", day, date.Day())
	assert.Equalf(t, month, date.Month(), "Expected matching month, wanted %d, got %d", month, date.Month())
	assert.Equalf(t, year, date.Year(), "Expected matching year, wanted %d, got %d", year, date.Year())
}

func Test_CreateTask(t *testing.T) {
	setup(t)

	task, _ := tm.CreateTask("build app", "create the task manager app")

	assert.Equal(t, "build app", task.Name)
	assert.Equal(t, "create the task manager app", task.Description)
	assert.False(t, task.Complete)
}

func Test_CreateTaskThenListTasks(t *testing.T) {
	setup(t)

	task, _ := tm.CreateTask("build app", "create the task manager app")

	taskSlice := tm.ListTasks()

	assert.Equal(t, 1, len(taskSlice), "Expected to find one task")
	assert.Contains(t, taskSlice, *task)
}

func Test_CreateTaskThenLoadTask(t *testing.T) {
	setup(t)

	tm.CreateTask("build app", "create the task manager app")

	task, _ := tm.GetTask("build app")

	assert.NotNil(t, task)
	assert.Equal(t, "build app", task.Name)
	assert.Equal(t, "create the task manager app", task.Description)
}

func Test_CompleteTask(t *testing.T) {
	setup(t)
	task, _ := tm.CreateTask("build app", "create the task manager app")

	task = tm.Complete(*task)

	assert.True(t, task.Complete)
	assertEqualCompletionDate(t, task.CompletionDate)
}

func Test_RemoveTask(t *testing.T) {
	setup(t)
	tm.CreateTask("build app", "create the task manager app")

	tm.Remove("build app")

	_, err := tm.GetTask("build app")

	assert.Equal(t, "Unable to locate a task with name build app", err.Error())
}

func Test_ListCompletedTasks(t *testing.T) {
	tm = TaskManager{
		repo: &MemRepository{
			tasks: []Task{
				{
					Name:           "Task 1",
					Description:    "Task Description",
					Complete:       true,
					CompletionDate: time.Now(),
				},
				{
					Name:           "Task 2",
					Description:    "Task Description",
					Complete:       true,
					CompletionDate: time.Now(),
				},
				{
					Name:           "Task 3",
					Description:    "Task Description",
					Complete:       true,
					CompletionDate: time.Now().Add(time.Duration(-8) * time.Hour),
				},
				{
					Name:           "Task 4",
					Description:    "Task Description",
					Complete:       true,
					CompletionDate: time.Now().Add(time.Duration(-15) * time.Hour),
				},
				{
					Name:           "Task 5",
					Description:    "Task Description",
					Complete:       true,
					CompletionDate: time.Now().Add(time.Duration(-25) * time.Hour),
				},
			},
		},
	}

	tests := []struct {
		name string
		want int
		args int
	}{
		{
			name: "Can load tasks completed today",
			want: 4,
			args: 24,
		},
		{
			name: "Can load tasks completed in last 12hrs",
			want: 3,
			args: 12,
		},
		{
			name: "Can load tasks completed in last 1 hr",
			want: 2,
			args: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks, _ := tm.ListCompletedTasks(tt.args)
			assert.Equalf(t, tt.want, len(tasks), "Expected to find %d tasks, got %d", tt.want, len(tasks))
		})
	}

}

type MemRepository struct {
	tasks []Task
	id    uint64
}

func (m *MemRepository) Store(t Task) (Task, error) {
	newTask := Task{
		Name:        t.Name,
		Description: t.Description,
		Complete:    t.Complete,
	}
	m.tasks = append(m.tasks, newTask)
	return newTask, nil
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

func (m *MemRepository) Load(name string) (*Task, error) {
	for _, v := range m.tasks {
		if v.Name == name {
			return &Task{
				Name:        v.Name,
				Description: v.Description,
				Complete:    v.Complete,
			}, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Unable to locate a task with name %s", name))
}

func (m *MemRepository) Delete(task Task) error {
	idx := -1
	for i, v := range m.tasks {
		if v.Name == task.Name {
			idx = i
			break
		}
	}
	if idx != -1 {
		if len(m.tasks) == 1 {
			m.tasks = []Task{}
		} else {
			m.tasks[idx] = m.tasks[len(m.tasks)-1]
			m.tasks = m.tasks[:len(m.tasks)-1]
		}
	}
	return nil
}

func (m *MemRepository) nextId() uint64 {
	m.id++
	return m.id
}
