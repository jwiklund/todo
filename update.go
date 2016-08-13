package main

import (
	"strings"

	"github.com/jwiklund/todo/ext"
	"github.com/jwiklund/todo/todo"
)

// todo [-v][-r <repo>] update <id> [-a <key> [<value>]][<state>]
func updateCmd(r ext.Repo, opts map[string]interface{}) {
	state := ""
	if s := opts["<state>"]; s != nil {
		state = s.(string)
	}
	key := ""
	if s := opts["<key>"]; s != nil {
		key = s.(string)
	}
	value := ""
	if s := opts["<value>"]; s != nil {
		value = s.(string)
	}
	var message []string
	if m := opts["<message>"]; m != nil {
		message = m.([]string)
	}
	update(r, opts["<id>"].(string), message, state, key, value)
}

//  todo [-v][-r <repo>] do <id>
func doCmd(r ext.Repo, opts map[string]interface{}) {
	update(r, opts["<id>"].(string), nil, "doing", "", "")
}

//  todo [-v][-r <repo>] wait <id>
func waitCmd(r ext.Repo, opts map[string]interface{}) {
	update(r, opts["<id>"].(string), nil, "waiting", "", "")
}

//  todo [-v][-r <repo>] done <id>
func doneCmd(r ext.Repo, opts map[string]interface{}) {
	update(r, opts["<id>"].(string), nil, "done", "", "")
}

func update(r ext.Repo, id string, message []string, state, key, value string) {
	task, err := r.Get(id)
	if err != nil {
		mainLog.Error(err.Error())
		mainLog.Debugf("%+v", err)
		return
	}
	if message != nil {
		task.Message = strings.Join(message, " ")
	}
	if state != "" {
		if !todo.StateValid(state) {
			mainLog.Debug("Invalid state ", state)
		}
		task.State = todo.StateFrom(state)
	}
	if key != "" {
		if value == "" {
			delete(task.Attr, key)
		} else {
			task.Attr[key] = value
		}
	}
	err = r.Update(task)
	if err != nil {
		mainLog.Error(err.Error())
		mainLog.Debugf("%+v", err)
		return
	}
	list(r, false, "")
}
