package main

import (
	"strings"

	"github.com/jwiklund/todo/ext"
)

// todo [-v][-r <repo>] add [-a <key> <value>] <message>...
func addCmd(r ext.Repo, opts map[string]interface{}) {
	messages := opts["<message>"].([]string)
	message := strings.Join(messages, " ")
	_, err := r.Add(message, nil)
	if err != nil {
		mainLog.Error("Could not add task ", err.Error())
		mainLog.Debugf("%+v", err)
		return
	}
	list(r, false, "")
}
