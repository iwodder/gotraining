package tasks

import (
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"testing"
)

var repo *BoltRepository

func setupTest(t *testing.T) func(t *testing.T) {
	if r, err := NewRepository("test.db"); err != nil {
		t.Fatal("Unable to open test.db ", err)
	} else {
		repo = r
	}
	return func(t *testing.T) {
		if err := os.Remove("test.db"); err != nil {
			t.Fatal("Unable to remove test.db ", err)
		}
	}
}

func TestBoltRepository_LoadAll(t *testing.T) {
	teardown := setupTest(t)
	defer teardown(t)

	repo.Store(Task{Name: "Task 1", Description: "Description 1"})
	repo.Store(Task{Name: "Task 2", Description: "Description 2"})

	tests := []struct {
		name string
		want []Task
	}{
		{
			"Can load a list of tasks",
			[]Task{
				{Name: "Task 1", Description: "Description 1"},
				{Name: "Task 2", Description: "Description 2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := repo.LoadAll(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoltRepository_Store(t *testing.T) {
	teardown := setupTest(t)
	defer teardown(t)

	type args struct {
		t Task
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Can store a task",
			args: args{
				t: Task{
					Name:        "Task 1",
					Description: "A description",
					Complete:    false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.Store(tt.args.t)
			task, _ := repo.Load(tt.args.t.Name)
			assert.Equal(t, tt.args.t, *task)
		})
	}
}

func TestBoltRepository_Update(t *testing.T) {
	teardown := setupTest(t)
	defer teardown(t)

	repo.Store(Task{Name: "Task 2", Description: "Initial Description", Complete: false})

	type args struct {
		task Task
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"Can update a task",
			args{
				Task{
					Name:        "Task 2",
					Description: "NewRepository Description",
					Complete:    false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := repo.Update(tt.args.task)
			assert.Equal(t, tt.args.task, actual)
		})
	}
}

func TestBoltRepository_LoadTask(t *testing.T) {
	teardown := setupTest(t)
	defer teardown(t)

	repo.Store(Task{Name: "Task 1", Description: "Initial Description", Complete: false})

	type args struct {
		taskName string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"Can load a task",
			args{taskName: "Task 1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, _ := repo.Load(tt.args.taskName)
			if !assert.NotNil(t, actual) {
				t.FailNow()
			}
			assert.Equal(t, tt.args.taskName, actual.Name)
		})
	}
}

func TestBoltRepository_RemoveTask(t *testing.T) {
	teardown := setupTest(t)
	defer teardown(t)

	repo.Store(Task{Name: "Task 1", Description: "Initial Description", Complete: false})

	type args struct {
		taskName string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"Can remove a task",
			args{taskName: "Task 1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.Delete(Task{Name: tt.args.taskName})

			_, err := repo.Load(tt.args.taskName)

			assert.Equal(t, "No task with name Task 1", err.Error())
		})
	}
}
