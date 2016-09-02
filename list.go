package main

import (
	"os"

	"github.com/jwiklund/todo/todo"
	"github.com/jwiklund/todo/view"
)

//  todo [-av][-r <repo>]
//  todo [-av][-r <repo>] list
func listCmd(t view.Todo, opts map[string]interface{}) {
	all := opts["-a"].(bool)
	state, _ := opts["<state>"].(string)
	list(t, all, state)
}

func list(t view.Todo, all bool, state string) {
	filter := func(t todo.Task) bool {
		return true
	}
	if !all {
		oldFilter := filter
		filter = func(t todo.Task) bool {
			return oldFilter(t) && t.IsCurrent()
		}
	}
	if state != "" {
		oldFilter := filter
		verifiedState := todo.StateFrom(state)
		filter = func(t todo.Task) bool {
			return oldFilter(t) && t.State == verifiedState
		}
	}
	tasks, err := t.List(filter)
	if err != nil {
		mainLog.Error("Couldn't list tasks ", err.Error())
		mainLog.Debugf("%+v", err)
		return
	}
	renderList(tasks, os.Stdout)
}
