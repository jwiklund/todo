package internal

import "github.com/jwiklund/todo/todo"

// Change a diff between two todo tasks
type Change struct {
	Added    map[string]string
	Modified map[string]string
	Removed  []string
}

// Compare return a Change between two tasks
func Compare(original, modified todo.Task) Change {
	change := Change{
		map[string]string{},
		map[string]string{},
		[]string{},
	}
	if original.Message != modified.Message {
		if original.Message == "" {
			change.Added["message"] = modified.Message
		} else if modified.Message == "" {
			change.Removed = append(change.Removed, "message")
		} else {
			change.Modified["message"] = modified.Message
		}
	}
	for key, value := range modified.Attr {
		if old, ok := original.Attr[key]; ok {
			if old != value {
				change.Modified[key] = value
			}
		} else {
			change.Added[key] = value
		}
	}
	for key := range original.Attr {
		if _, ok := modified.Attr[key]; !ok {
			change.Removed = append(change.Removed, key)
		}
	}
	return change
}

// Apply change to task, modifies task
func (c Change) Apply(original *todo.Task) {
	for _, key := range c.Removed {
		if key == "message" {
			original.Message = ""
		} else {
			delete(original.Attr, key)
		}
	}
	for key, value := range c.Modified {
		if key == "message" {
			original.Message = value
		} else {
			original.Attr[key] = value
		}
	}
	for key, value := range c.Added {
		if key == "message" {
			original.Message = value
		} else {
			original.Attr[key] = value
		}
	}
}
