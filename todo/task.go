package todo

import "reflect"

// Task a todo task
type Task struct {
	ID      string
	State   State
	Message string
	Attr    map[string]string
}

func (t Task) String() string {
	return "(" + t.ID + ")\t" + t.State.String() + "\t" + t.Message
}

// Equal check if tasks are equal
func (t Task) Equal(t2 Task) bool {
	return reflect.DeepEqual(t, t2)
}

// IsCurrent return true if task is not waiting or archived
func (t Task) IsCurrent() bool {
	return t.State == "todo" || t.State == "doing"
}
