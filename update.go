package main

import (
	"strings"

	"github.com/jwiklund/todo/todo"
	"github.com/jwiklund/todo/view"
)

// todo [-v][-r <repo>] update <id> [-a <key> [<value>]][<state>]
func updateCmd(t view.Todo, opts map[string]interface{}) {
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
	update(t, opts["<id>"].(string), message, state, key, value)
}

//  todo [-v][-r <repo>] do <id>
func doCmd(t view.Todo, opts map[string]interface{}) {
	update(t, opts["<id>"].(string), nil, "doing", "", "")
}

//  todo [-v][-r <repo>] wait <id>
func waitCmd(t view.Todo, opts map[string]interface{}) {
	update(t, opts["<id>"].(string), nil, "waiting", "", "")
}

//  todo [-v][-r <repo>] done <id>
func doneCmd(t view.Todo, opts map[string]interface{}) {
	update(t, opts["<id>"].(string), nil, "done", "", "")
}

func update(t view.Todo, id string, message []string, state, key, value string) {
	task, err := t.Get(id)
	if err != nil {
		mainLog.Error(err.Error())
		mainLog.Debugf("%+v", err)
		return
	}
	if len(message) != 0 {
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
	err = t.Update(task)
	if err != nil {
		mainLog.Error(err.Error())
		mainLog.Debugf("%+v", err)
		return
	}
	list(t, false, "")
}
