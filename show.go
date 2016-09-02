package main

import (
	"os"

	"github.com/jwiklund/todo/view"
)

//  todo [-v][-r <repo>] show <id>
func showCmd(t view.Todo, opts map[string]interface{}) {
	task, err := t.Get(opts["<id>"].(string))
	if err != nil {
		mainLog.Error(err.Error())
		mainLog.Debugf("%+v", err)
	}
	renderOne(task, os.Stdout)
}
