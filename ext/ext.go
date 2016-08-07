package ext

import "github.com/jwiklund/todo/todo"

// External storage interface
type External interface {
	Handle(task todo.Task) (todo.Task, error)
}

var ext external

// Register an external extension
func Register(e External) {
	ext.Register(e)
}

// Registered an external for all registered externals
func Registered() External {
	return &ext
}

type external struct {
	externals []External
}

func (ext *external) Register(e External) {
	ext.externals = append(ext.externals, e)
}

func (ext *external) Handle(task todo.Task) (todo.Task, error) {
	for _, external := range ext.externals {
		mod, err := external.Handle(task)
		if err != nil {
			return mod, err
		}
		task = mod
	}
	return task, nil
}
