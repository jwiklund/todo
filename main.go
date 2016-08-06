package main

import (
	"bytes"
	"fmt"

	log "github.com/Sirupsen/logrus"
	opt "github.com/docopt/docopt-go"
	"github.com/jwiklund/todo/todo"
	_ "github.com/pkg/errors"
)

var usage = `Todo list.

Usage:
  todo -h
  todo [-av][-r <repo>] list
  todo [-av][-r <repo>] add <message...>
  todo [-av][-r <repo>] update <id> <state>
    
Options:
  -a          include all tasks [default false]
  -v          be verbose (debug) [default false]
  -r <repo>   custom repo [default sqlite://~/.todo.db]
`
var mainLog = log.WithField("comp", "main")

func main() {
	opts, err := opt.Parse(usage, nil, true, "1.0", false)
	if err != nil {
		log.Fatal(err)
		return
	}
	if opts["-v"].(bool) {
		log.SetLevel(log.DebugLevel)
	}
	mainLog.Debug("Args ", sortOpts(opts))

	repo := repo(opts)
	if repo != nil {
		defer repo.Close()
		cmd(repo, opts)
	}
}

func repo(opts map[string]interface{}) todo.Repo {
	path, _ := opts["-r"].(string)
	if path == "" {
		path = "sqlite://~/.todo.db"
	}
	repo, err := todo.RepoFromPath(path)
	if err != nil {
		mainLog.Error("Invalid repository path ", path, " ", err.Error())
		mainLog.Debugf("%+v", err)
		return nil
	}
	return repo
}

func cmd(r todo.Repo, opts map[string]interface{}) {
	if opts["list"] != nil {
		all := opts["-a"].(bool)
		tasks, err := r.List()
		if err != nil {
			mainLog.Error("Couldn't list tasks ", err.Error())
			mainLog.Debugf("%+v", err)
			return
		}
		for _, task := range tasks {
			if all || task.IsCurrent() {
				fmt.Println(task)
			}
		}
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
