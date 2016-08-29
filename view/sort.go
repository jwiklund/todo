package view

import (
	"strconv"

	"github.com/jwiklund/todo/todo"
)

type sorter []todo.Task

func (s sorter) Len() int {
	return len(s)
}

func (s sorter) Swap(i, j int) {
	t := s[i]
	s[i] = s[j]
	s[j] = t
}

func (s sorter) Less(i, j int) bool {
	p1 := prio(s[i])
	p2 := prio(s[j])
	if p1 < p2 {
		return true
	}
	if p2 < p1 {
		return false
	}

	return id(s[i]) < id(s[j])
}

func prio(t todo.Task) int {
	if p, ok := t.Attr["prio"]; ok {
		if prio, err := strconv.Atoi(p); err == nil {
			return prio
		}
	}
	return 1000
}

func id(t todo.Task) int {
	if i, err := strconv.Atoi(t.ID); err == nil {
		return i
	}
	return 0
}
