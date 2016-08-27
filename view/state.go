package view

import (
	"strconv"

	"github.com/jwiklund/todo/todo"
	"github.com/pkg/errors"
)

// State view state
type State struct {
	// Mapping from view ids to task id
	Mapping map[string]string
}

// ToDB relative to absolute id
func (s *State) ToDB(id string) (string, error) {
	if mapped, ok := s.Mapping[id]; ok {
		return mapped, nil
	}
	return id, errors.New("Relative ID not found")
}

// Remapp remap tasks and rewrite state mapping
func (s *State) Remapp(tasks []todo.Task) ([]todo.Task, error) {
	r := make([]todo.Task, len(tasks))
	s.Mapping = map[string]string{}
	for i, task := range tasks {
		r[i] = task
		relID := strconv.Itoa(i)
		s.Mapping[relID] = r[i].ID
		r[i].ID = relID
	}
	return r, nil
}
