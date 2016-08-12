package ext

import "github.com/jwiklund/todo/todo"

// Repo an external repo
type Repo interface {
	todo.Repo

	Sync(name string) error
	SyncAll() error
}

// ExternalRepo wrap a repo with external handling
func ExternalRepo(repo todo.RepoBegin, config []ExternalConfig) (Repo, error) {
	external, err := configure(config)
	if err != nil {
		return nil, err
	}
	return &extRepo{repo, external}, nil
}

type extRepo struct {
	repo todo.RepoBegin
	ext  external
}

func (r *extRepo) Sync(name string) error {
	return r.ext.Sync(r.repo, name)
}

func (r *extRepo) SyncAll() error {
	return r.ext.SyncAll(r.repo)
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

func (r *extRepo) AddWithAttr(message string, attr map[string]string) (todo.Task, error) {
	task, err := r.repo.AddWithAttr(message, attr)
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
