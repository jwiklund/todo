package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/jwiklund/todo/ext"
	_ "github.com/jwiklund/todo/ext/text"
	"github.com/jwiklund/todo/todo"
	"github.com/jwiklund/todo/util"

	"github.com/BurntSushi/toml"
	"github.com/Sirupsen/logrus"
	opt "github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var usage = `Todo list.

Usage:
  todo -h
  todo [(-r <repo>) (-c <cfg>) -av]
  todo [(-r <repo>) (-c <cfg>) -av] list [<state>]
  todo [(-r <repo>) (-c <cfg>) -v] add [(-a <key> <value>)] <message>...
  todo [(-r <repo>) (-c <cfg>) -v] update <id> [(-a <key> <value>) (-s <state>) (-m <message>...)]
  todo [(-r <repo>) (-c <cfg>) -v] sync [<external>]
  todo [(-r <repo>) (-c <cfg>) -v] show <id>
  todo [(-r <repo>) (-c <cfg>) -v] do <id>
  todo [(-r <repo>) (-c <cfg>) -v] wait <id>
  todo [(-r <repo>) (-c <cfg>) -v] done <id>
    
Options:
  -a          include all tasks [default false]
  -v          be verbose (debug) [default false]
  -c <cfg>    config [default ~/.todo.conf]
  -r <repo>   custom repo [default sqlite://~/.todo.db]
`
var mainLog = logrus.WithField("comp", "main")

var cmds = map[string]func(ext.Repo, map[string]interface{}){
	"list":   listCmd,
	"add":    addCmd,
	"update": updateCmd,
	"sync":   syncCmd,
	"show":   showCmd,
	"do":     doCmd,
	"wait":   waitCmd,
	"done":   doneCmd,
}

type config struct {
	External []ext.ExternalConfig
}

func main() {
	opts, err := opt.Parse(usage, nil, true, "1.0", false)
	if err != nil {
		mainLog.Fatal(err)
		return
	}
	if opts["-v"].(bool) {
		logrus.SetLevel(logrus.DebugLevel)
	}
	mainLog.Debug("Args ", sortOpts(opts))

	repo := repo(opts)
	if repo == nil {
		return
	}

	config, err := readConfig(opts["<cfg>"])
	if err != nil {
		mainLog.Error(err.Error())
		mainLog.Debug("%+v", err)
		return
	}
	extRepo, err := ext.ExternalRepo(repo, config.External)
	if err != nil {
		mainLog.Error(err.Error())
		mainLog.Debug("%+v", err)
		return
	}
	cmd(extRepo, opts)
	mainLog.Debugf("Close returned %v", extRepo.Close())
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

func repo(opts map[string]interface{}) todo.RepoBegin {
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

func readConfig(path interface{}) (config, error) {
	p := "~/.todo.conf"
	c := config{}
	if path != nil {
		p = path.(string)
	}
	p = util.Expand(p)
	stat, err := os.Stat(p)
	if err != nil {
		mainLog.Debug("Config file not found, using defaults", err)
		return c, nil
	}
	if stat.IsDir() {
		return c, errors.Wrap(err, "Config is a directory")
	}

	if _, err := toml.DecodeFile(p, &c); err != nil {
		return c, errors.Wrap(err, "Could not read config")
	}
	return c, nil
}

func cmd(r ext.Repo, opts map[string]interface{}) {
	for key, cmd := range cmds {
		if opts[key].(bool) {
			cmd(r, opts)
			return
		}
	}
	listCmd(r, opts)
}

//  todo [-av][-r <repo>]
//  todo [-av][-r <repo>] list
func listCmd(r ext.Repo, opts map[string]interface{}) {
	all := opts["-a"].(bool)
	state, _ := opts["<state>"].(string)
	list(r, all, state)
}

func list(r ext.Repo, all bool, state string) {
	tasks, err := r.List()
	if err != nil {
		mainLog.Error("Couldn't list tasks ", err.Error())
		mainLog.Debugf("%+v", err)
		return
	}
	filter := func(t todo.Task) bool {
		return true
	}
	if !all {
		oldFilter := filter
		filter = func(t todo.Task) bool {
			return oldFilter(t) && t.IsCurrent()
		}
	}
	if state != "" {
		oldFilter := filter
		verifiedState := todo.StateFrom(state)
		filter = func(t todo.Task) bool {
			return oldFilter(t) && t.State == verifiedState
		}
	}
	w := tabwriter.NewWriter(os.Stdout, 6, 8, 1, ' ', 0)
	for _, task := range tasks {
		if filter(task) {
			w.Write([]byte(task.String()))
			w.Write([]byte{'\n'})
		}
	}
	w.Flush()
}

// todo [-v][-r <repo>] add [-a <key> <value>] <message>...
func addCmd(r ext.Repo, opts map[string]interface{}) {
	messages := opts["<message>"].([]string)
	message := strings.Join(messages, " ")
	_, err := r.Add(message)
	if err != nil {
		mainLog.Error("Could not add task ", err.Error())
		mainLog.Debugf("%+v", err)
		return
	}
	list(r, false, "")
}

//  todo [-v][-r <repo>] show <id>
func showCmd(r ext.Repo, opts map[string]interface{}) {
	task, err := r.Get(opts["<id>"].(string))
	if err != nil {
		mainLog.Error(err.Error())
		mainLog.Debugf("%+v", err)
	}
	w := tabwriter.NewWriter(os.Stdout, 6, 8, 2, ' ', 0)
	fmt.Fprintf(w, "(%s)\t%s\t%s\n", task.ID, task.State.String(), task.Message)
	for key, value := range task.Attr {
		fmt.Fprintf(w, "\t%s\t%s\n", key, value)
	}
	w.Flush()
}

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
