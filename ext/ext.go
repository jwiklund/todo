package ext

import (
	"errors"

	"github.com/jwiklund/todo/todo"
)

// External storage interface
type External interface {
	Handle(task todo.Task) (todo.Task, error)
	Sync(r todo.Repo) error
	Close() error
}

// ExternalConfig external storage config
type ExternalConfig struct {
	Type string
	ID   string
	URI  string
}

// ExternalFactory factory for externals
type ExternalFactory func(ExternalConfig) (External, error)

var externalFuncs map[string]ExternalFactory

// Register a named external factory
func Register(name string, factory ExternalFactory) {
	if externalFuncs == nil {
		externalFuncs = map[string]ExternalFactory{
			name: factory,
		}
	} else {
		externalFuncs[name] = factory
	}
}

func configure(configs []ExternalConfig) (external, error) {
	ext := external{}

	for _, config := range configs {
		factory, ok := externalFuncs[config.Type]
		if !ok {
			return ext, errors.New("External type not found " + config.Type)
		}
		e, err := factory(config)
		if err != nil {
			return ext, err
		}
		if ext.externals == nil {
			ext.externals = map[string]External{config.ID: e}
		} else {
			ext.externals[config.ID] = e
		}
	}

	return ext, nil
}

type external struct {
	externals map[string]External
}

func (ext external) Externals() map[string]External {
	return ext.externals
}

func (ext external) Close() error {
	var e error
	for _, external := range ext.externals {
		err := external.Close()
		if err != nil && e == nil {
			e = err
		}
	}
	return e
}

func (ext external) Handle(task todo.Task) (todo.Task, error) {
	for _, external := range ext.externals {
		mod, err := external.Handle(task)
		if err != nil {
			return mod, err
		}
		task = mod
	}
	return task, nil
}
