package ext

import "github.com/jwiklund/todo/todo"

// Repo wrap a repo with external handling
func Repo(repo todo.Repo) todo.Repo {
	return &extRepo{repo}
}

type extRepo struct {
	repo todo.Repo
}

func (r *extRepo) Close() error {
	return r.repo.Close()
}

func (r *extRepo) List() ([]todo.Task, error) {
	return r.repo.List()
}

func (r *extRepo) Add(message string) (todo.Task, error) {
	task, err := r.repo.Add(message)
	if err != nil {
		return task, err
	}
	upd, err := ext.Handle(task)
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
	mod, err := ext.Handle(task)
	if err != nil {
		return err
	}
	return r.repo.Update(mod)
}
