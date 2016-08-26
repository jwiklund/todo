package ext

import (
	"github.com/jwiklund/todo/todo"
	"github.com/pkg/errors"
)

// External (single) storage interface
type External interface {
	Handle(task todo.Task) (todo.Task, error)
	Sync(r todo.RepoBegin, dryRun bool) error
	Close() error
}

// Externals storage interface
type Externals interface {
	Handle(task todo.Task) (todo.Task, error)
	SyncAll(r todo.RepoBegin, dryRun bool) error
	Sync(r todo.RepoBegin, name string, dryRun bool) error
	Close() error
}

// ExternalConfig external storage config
type ExternalConfig struct {
	Type  string
	ID    string
	URI   string
	Extra map[string]string
}

// ExternalFactory factory for externals
type ExternalFactory func(ExternalConfig) (External, error)

var externalFuncs map[string]ExternalFactory

// Register a named external factory
func Register(name string, factory ExternalFactory) {
	if externalFuncs == nil {
		externalFuncs = map[string]ExternalFactory{name: factory}
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

func (ext external) Sync(r todo.RepoBegin, name string, dryRun bool) error {
	if ext, ok := ext.externals[name]; ok {
		return ext.Sync(r, dryRun)
	}
	return errors.Errorf("external %s does not exist", name)
}

func (ext external) SyncAll(r todo.RepoBegin, dryRun bool) error {
	for id, ext := range ext.externals {
		err := ext.Sync(r, dryRun)
		if err != nil {
			return errors.Wrapf(err, "Failed %s", id)
		}
	}

	return nil
}
