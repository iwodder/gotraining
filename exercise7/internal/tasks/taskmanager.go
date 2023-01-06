package tasks

type Task struct {
	ID          uint64
	Name        string
	Description string
	Completed   bool
}

type TaskManager struct {
	repo TaskRepository
}

type TaskRepository interface {
	Store(t Task) Task
	LoadAll() []Task
	Update(task Task) Task
	LoadTask(name string) (*Task, error)
}

func NewTaskManager(t TaskRepository) *TaskManager {
	return &TaskManager{
		repo: t,
	}
}

func (tm *TaskManager) CreateTask(name string, description string) *Task {
	t := tm.repo.Store(Task{
		Name:        name,
		Description: description,
	})
	return &t
}

func (tm *TaskManager) GetTask(name string) (*Task, error) {
	return tm.repo.LoadTask(name)
}

func (tm *TaskManager) ListTasks() []Task {
	return tm.repo.LoadAll()
}

func (tm *TaskManager) Complete(task Task) *Task {
	task.Completed = true
	t := tm.repo.Update(task)
	return &t
}
