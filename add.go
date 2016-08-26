package main

import (
	"strings"

	"github.com/jwiklund/todo/view"
)

// todo [-v][-r <repo>] add [-a <key> <value>] <message>...
func addCmd(t view.Todo, opts map[string]interface{}) {
	messages := opts["<message>"].([]string)
	message := strings.Join(messages, " ")
	_, err := t.Add(message, nil)
	if err != nil {
		mainLog.Error("Could not add task ", err.Error())
		mainLog.Debugf("%+v", err)
		return
	}
	list(t, false, "")
}
