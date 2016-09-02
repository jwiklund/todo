package main

import (
	"fmt"

	"github.com/jwiklund/todo/view"
)

func externalCmd(t view.Todo, opts map[string]interface{}) {
	task, err := t.Get(opts["<id>"].(string))
	if err != nil {
		mainLog.Info(err.Error())
		mainLog.Debugf("%+v", err)
		return
	}
	if ext, _ := opts["<external>"]; ext != nil {
		if task.External() != "" {
			mainLog.Info("Task already has external")
			return
		}
		update(t, opts["<id>"].(string), nil, "", "external", ext.(string))
	} else {
		ext := task.External()
		if ext == "" {
			ext = "<none>"
		}
		fmt.Printf("(%s)   %s  %s\n", task.ID, ext, task.ExternalID())
	}
}
