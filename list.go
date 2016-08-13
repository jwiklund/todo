package main

import (
	"os"
	"text/tabwriter"

	"github.com/jwiklund/todo/ext"
	"github.com/jwiklund/todo/todo"
)

//  todo [-av][-r <repo>]
//  todo [-av][-r <repo>] list
func listCmd(r ext.Repo, opts map[string]interface{}) {
	all := opts["-a"].(bool)
	state, _ := opts["<state>"].(string)
	list(r, all, state)
}

func list(r ext.Repo, all bool, state string) {
	tasks, err := r.List()
	if err != nil {
		mainLog.Error("Couldn't list tasks ", err.Error())
		mainLog.Debugf("%+v", err)
		return
	}
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
	w := tabwriter.NewWriter(os.Stdout, 6, 8, 1, ' ', 0)
	for _, task := range tasks {
		if filter(task) {
			w.Write([]byte(task.String()))
			w.Write([]byte{'\n'})
		}
	}
	w.Flush()
}
