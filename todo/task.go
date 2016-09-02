package todo

import (
	"reflect"
	"strconv"
)

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

// Prio get the prio of this task
func (t Task) Prio() int {
	prio, _ := t.Attr["prio"]
	if prio == "" {
		return 1000
	}
	prioInt, err := strconv.Atoi(prio)
	if err != nil {
		todoLog.Debugf("Invalid prio %s", prio)
		prioInt = 1000
	}
	return prioInt
}

// SetPrio set prio of this task
func (t Task) SetPrio(prio int) {
	t.Attr["prio"] = strconv.Itoa(prio)
}
