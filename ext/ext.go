package ext

import (
	"errors"

	"github.com/jwiklund/todo/todo"
)

// External storage interface
type External interface {
	Handle(task todo.Task) (todo.Task, error)
	Close() error
}

// ExternalConfig external storage config
type ExternalConfig struct {
	Type string
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

func Configure(configs []ExternalConfig) (External, error) {
	ext := external{}

	for _, config := range configs {
		factory, ok := externalFuncs[config.Type]
		if !ok {
			return nil, errors.New("External type not found " + config.Type)
		}
		e, err := factory(config)
		if err != nil {
			return nil, err
		}
		ext.Register(e)
	}

	return &ext, nil
}

type external struct {
	externals []External
}

func (ext *external) Register(e External) {
	ext.externals = append(ext.externals, e)
}

func (ext *external) Close() error {
	var e error
	for _, external := range ext.externals {
		err := external.Close()
		if err != nil && e == nil {
			e = err
		}
	}
	return e
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
