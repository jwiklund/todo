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

// External get external repo identifier, if present (else "")
func (t Task) External() string {
	ext, _ := t.Attr["external"]
	return ext
}

// ExternalID get id in external repo, if present (else "")
func (t Task) ExternalID() string {
	ext := t.External()
	extID, _ := t.Attr[ext+".id"]
	return extID
}

// SetExternalID set external id if external is present
func (t Task) SetExternalID(externalID string) {
	if t.External() == "" {
		return
	}
	t.Attr[t.External()+".id"] = externalID
}
