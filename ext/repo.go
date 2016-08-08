package ext

import (
	"github.com/jwiklund/todo/todo"
)

// Repo wrap a repo with external handling
func Repo(repo todo.Repo, config []ExternalConfig) (todo.Repo, error) {
	external, err := Configure(config)
	if err != nil {
		return nil, err
	}
	return &extRepo{repo, external}, nil
}

type extRepo struct {
	repo todo.Repo
	ext  External
}

func (r *extRepo) Close() error {
	e1 := r.repo.Close()
	e2 := r.ext.Close()
	if e1 != nil {
		return e1
	}
	if e2 != nil {
		return e2
	}
	return nil
}

func (r *extRepo) List() ([]todo.Task, error) {
	return r.repo.List()
}

func (r *extRepo) Add(message string) (todo.Task, error) {
	task, err := r.repo.Add(message)
	if err != nil {
		return task, err
	}
	upd, err := r.ext.Handle(task)
	if err != nil {
		return upd, err
	}
	if !task.Equal(upd) {
		return upd, r.repo.Update(upd)
	}
	return upd, nil
}

func (r *extRepo) Get(id string) (todo.Task, error) {
	return r.repo.Get(id)
}

func (r *extRepo) Update(task todo.Task) error {
	mod, err := r.ext.Handle(task)
	if err != nil {
		return err
	}
	return r.repo.Update(mod)
}
