package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/Sirupsen/logrus"
	opt "github.com/docopt/docopt-go"
	"github.com/jwiklund/todo/ext"
	"github.com/jwiklund/todo/todo"
)

var usage = `Todo list.

Usage:
  todo -h
  todo [-av][-r <repo>]
  todo [-av][-r <repo>] list
  todo [-v][-r <repo>] add [-a <key> <value>] <message>...
  todo [-v][-r <repo>] update <id> [-a <key> [<value>]][<state>]
  todo [-v][-r <repo>] show <id>
  todo [-v][-r <repo>] do <id>
  todo [-v][-r <repo>] wait <id>
  todo [-v][-r <repo>] done <id>
    
Options:
  -a          include all tasks [default false]
  -v          be verbose (debug) [default false]
  -r <repo>   custom repo [default sqlite://~/.todo.db]
`
var log = logrus.WithField("comp", "main")

var cmds = map[string]func(todo.Repo, map[string]interface{}){
	"list":   listCmd,
	"add":    addCmd,
	"update": updateCmd,
	"show":   showCmd,
	"do":     doCmd,
	"wait":   waitCmd,
	"done":   doneCmd,
}

func main() {
	opts, err := opt.Parse(usage, nil, true, "1.0", false)
	if err != nil {
		log.Fatal(err)
		return
	}
	if opts["-v"].(bool) {
		logrus.SetLevel(logrus.DebugLevel)
	}
	log.Debug("Args ", sortOpts(opts))

	repo := repo(opts)
	if repo != nil {
		extRepo := ext.Repo(repo)
		cmd(extRepo, opts)
		log.Debugf("Close returned %v", extRepo.Close())
	}
}

func sortOpts(opts map[string]interface{}) string {
	res := bytes.Buffer{}
	var keys []string
	for key := range opts {
		keys = append(keys, key)
	}
	for i, key := range keys {
		if i != 0 {
			res.WriteString(" ")
		}
		res.WriteString(fmt.Sprintf("%s=%v", key, opts[key]))
	}
	return res.String()
}

func repo(opts map[string]interface{}) todo.Repo {
	path, _ := opts["-r"].(string)
	if path == "" {
		path = "sqlite://~/.todo.db"
	}
	repo, err := todo.RepoFromPath(path)
	if err != nil {
		log.Error("Invalid repository path ", path, " ", err.Error())
		log.Debugf("%+v", err)
		return nil
	}
	return repo
}

func cmd(r todo.Repo, opts map[string]interface{}) {
	for key, cmd := range cmds {
		if opts[key].(bool) {
			cmd(r, opts)
			return
		}
	}
	listCmd(r, opts)
}

// todo [-v][-r <repo>] add [-a <key> <value>] <message>...
func addCmd(r todo.Repo, opts map[string]interface{}) {
	messages := opts["<message>"].([]string)
	message := strings.Join(messages, " ")
	task, err := r.Add(message)
	if err != nil {
		log.Error("Could not add task ", err.Error())
		log.Debugf("%+v", err)
		return
	}
	fmt.Println("Created ", task)
}

//  todo [-av][-r <repo>]
//  todo [-av][-r <repo>] list
func listCmd(r todo.Repo, opts map[string]interface{}) {
	all := opts["-a"].(bool)
	tasks, err := r.List()
	if err != nil {
		log.Error("Couldn't list tasks ", err.Error())
		log.Debugf("%+v", err)
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 6, 8, 1, ' ', 0)
	for _, task := range tasks {
		if all || task.IsCurrent() {
			w.Write([]byte(task.String()))
			w.Write([]byte{'\n'})
		}
	}
	w.Flush()
}

//  todo [-v][-r <repo>] show <id>
func showCmd(r todo.Repo, opts map[string]interface{}) {
	task, err := r.Get(opts["<id>"].(string))
	if err != nil {
		log.Error(err.Error())
		log.Debugf("%+v", err)
	}
	w := tabwriter.NewWriter(os.Stdout, 6, 8, 2, ' ', 0)
	fmt.Fprintf(w, "(%s)\t%s\t%s\n", task.ID, task.State.String(), task.Message)
	for key, value := range task.Attr {
		fmt.Fprintf(w, "\t%s\t%s\n", key, value)
	}
	w.Flush()
}

// todo [-v][-r <repo>] update <id> [-a <key> [<value>]][<state>]
func updateCmd(r todo.Repo, opts map[string]interface{}) {
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
	update(r, opts["<id>"].(string), state, key, value)
}

//  todo [-v][-r <repo>] do <id>
func doCmd(r todo.Repo, opts map[string]interface{}) {
	update(r, opts["<id>"].(string), "doing", "", "")
}

//  todo [-v][-r <repo>] wait <id>
func waitCmd(r todo.Repo, opts map[string]interface{}) {
	update(r, opts["<id>"].(string), "waiting", "", "")
}

//  todo [-v][-r <repo>] done <id>
func doneCmd(r todo.Repo, opts map[string]interface{}) {
	update(r, opts["<id>"].(string), "done", "", "")
}

func update(r todo.Repo, id, state, key, value string) {
	task, err := r.Get(id)
	if err != nil {
		log.Error(err.Error())
		log.Debugf("%+v", err)
		return
	}
	if state != "" {
		if !todo.StateValid(state) {
			log.Debug("Invalid state ", state)
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
		log.Error(err.Error())
		log.Debugf("%+v", err)
		return
	}
}
