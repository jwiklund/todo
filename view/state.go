package view

import "github.com/pkg/errors"

// State view state
type State struct {
	// Mapping from view ids to task id
	Mapping map[string]string
}

// ToDB relative to absolute id
func (s State) ToDB(id string) (string, error) {
	if mapped, ok := s.Mapping[id]; ok {
		return mapped, nil
	}
	return id, errors.New("Relative ID not found")
}
