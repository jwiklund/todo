package jira

import "github.com/jwiklund/todo/todo"

// TODO handle unknown states (do not change state unless explicitly requested)
func jiraState(jiraStatus string) todo.State {
	switch jiraStatus {
	case "To Do":
		return todo.StateTodo
	case "In Progress":
		return todo.StateDoing
	case "Done":
		return todo.StateDone
	default:
		return todo.StateTodo
	}
}

func jiraStatus(state todo.State) string {
	switch state {
	case todo.StateTodo:
		return "To Do"
	case todo.StateDoing:
		return "In Progress"
	case todo.StateWaiting:
		return "To Do"
	case todo.StateDone:
		return "Done"
	default:
		return "To Do"
	}
}
