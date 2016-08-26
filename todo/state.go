package todo

// State valid state
type State string

var (
	// StateTodo "todo"
	StateTodo = State("todo")
	// StateWaiting "waiting"
	StateWaiting = State("waiting")
	// StateDoing "doing"
	StateDoing = State("doing")
	// StateDone "done"
	StateDone = State("done")
	// States all states
	States = []State{StateTodo, StateWaiting, StateDoing, StateDone}
)

func (s State) String() string {
	return string(s)
}

// StateValid check if state is a valid state
func StateValid(state string) bool {
	return state == StateFrom(state).String()
}

// StateFrom returns a valid state from a given string (or todo)
func StateFrom(state string) State {
	switch state {
	case "todo":
		return State(state)
	case "doing":
		return State(state)
	case "waiting":
		return State(state)
	case "done":
		return State(state)
	}
	return "todo"
}
