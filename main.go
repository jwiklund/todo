package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/Sirupsen/logrus"
	opt "github.com/docopt/docopt-go"
	"github.com/jwiklund/todo/todo"
)

var usage = `Todo list.

Usage:
  todo -h
  todo [-av][-r <repo>]
  todo [-av][-r <repo>] list
  todo [-v][-r <repo>] add <message>...
  todo [-v][-r <repo>] update <id> <state>
  todo [-v][-r <repo>] do <id>
  todo [-v][-r <repo>] wait <id>
  todo [-v][-r <repo>] done <id>
    
Options:
  -a          include all tasks [default false]
  -v          be verbose (debug) [default false]
  -r <repo>   custom repo [default sqlite://~/.todo.db]
`
var log = logrus.WithField("comp", "main")

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
		defer repo.Close()
		cmd(repo, opts)
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
	if opts["add"].(bool) {
		messages := opts["<message>"].([]string)
		message := strings.Join(messages, " ")
		task, err := r.Add(message)
		if err != nil {
			log.Error("Could not add task ", err.Error())
			log.Debugf("%+v", err)
			return
		}
		fmt.Println("Created ", task)
		return
	}
	if opts["update"].(bool) {
		update(r, opts["<id>"].(string), opts["<state>"].(string))
		return
	}
	if opts["do"].(bool) {
		update(r, opts["<id>"].(string), "doing")
		return
	}
	if opts["wait"].(bool) {
		update(r, opts["<id>"].(string), "waiting")
		return
	}
	if opts["done"].(bool) {
		update(r, opts["<id>"].(string), "done")
		return
	}
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

func update(r todo.Repo, id, state string) {
	if !todo.StateValid(state) {
		log.Debug("Invalid state ", state)
	}
	s := todo.StateFrom(state)
	task, err := r.Get(id)
	if err != nil {
		log.Error("Could not find task ", err.Error())
		log.Debugf("%+v", err)
		return
	}
	task.State = s
	err = r.Update(task)
	if err != nil {
		log.Error("Could not update task ", err.Error())
		log.Debugf("%+v", err)
		return
	}
}
