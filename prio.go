package main

import (
	"fmt"

	"github.com/jwiklund/todo/view"
)

// todo [-v][-r <repo>] prio <id> [<prio>]
func prioCmd(t view.Todo, opts map[string]interface{}) {
	if prio, _ := opts["<prio>"]; prio != nil {
		update(t, opts["<id>"].(string), nil, "", "prio", prio.(string))
	} else {
		task, err := t.Get(opts["<id>"].(string))
		if err != nil {
			mainLog.Info(err.Error())
			mainLog.Debugf("%+v", err)
			return
		}
		fmt.Printf("%d\n", task.Prio())
	}
}
