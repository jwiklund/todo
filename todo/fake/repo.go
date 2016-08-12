package fake

import (
	"errors"
	"strconv"

	"github.com/jwiklund/todo/todo"
)

// New return a fake repo
func New() *Fake {
	return &Fake{}
}

// Fake a fake repo
type Fake struct {
	todos []todo.Task
}

// Close noop
func (r *Fake) Close() error {
	return nil
}

// Begin noop
func (r *Fake) Begin() (todo.RepoCommit, error) {
	return r, nil
}

// Commit noop
func (r *Fake) Commit() error {
	return nil
}

// List return todos
func (r *Fake) List() ([]todo.Task, error) {
	return r.todos, nil
}

// MustList return todos
func (r *Fake) MustList() []todo.Task {
	return r.todos
}

// Get return task
func (r *Fake) Get(id string) (todo.Task, error) {
	for _, t := range r.todos {
		if t.ID == id {
			return t, nil
		}
	}
	return todo.Task{}, errors.New("not found")
}

// MustGet return task or empty
func (r *Fake) MustGet(id string) todo.Task {
	for _, t := range r.todos {
		if t.ID == id {
			return t
		}
	}
	return todo.Task{}
}

// Add create task
func (r *Fake) Add(message string) (todo.Task, error) {
	return r.AddWithAttr(message, make(map[string]string))
}

// AddWithAttr create task
func (r *Fake) AddWithAttr(message string, attr map[string]string) (todo.Task, error) {
	task := todo.Task{
		ID:      strconv.Itoa(len(r.todos)),
		Message: message,
		Attr:    attr,
		State:   todo.StateFrom("todo"),
	}
	r.todos = append(r.todos, task)
	return task, nil
}

// Update update task
func (r *Fake) Update(newTask todo.Task) error {
	for i, task := range r.todos {
		if task.ID == newTask.ID {
			r.todos[i] = newTask
			return nil
		}
	}
	return errors.New("task not found")
}
