package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/BurntSushi/toml"
	"github.com/Sirupsen/logrus"
	opt "github.com/docopt/docopt-go"
	"github.com/jwiklund/todo/ext"
	_ "github.com/jwiklund/todo/ext/jira"
	_ "github.com/jwiklund/todo/ext/text"
	"github.com/jwiklund/todo/todo"
	"github.com/jwiklund/todo/util"
	"github.com/pkg/errors"
)

var usage = `Todo list.

Usage:
  todo -h
  todo [(-c <cfg>) -va]
  todo [(-c <cfg>) -va] list [<state>]
  todo [(-c <cfg>) -v] add [(-a <key> <value>)] <message>...
  todo [(-c <cfg>) -v] update <id> [(-a <key> <value>) (-s <state>) (-m <message>...)]
  todo [(-c <cfg>) -vd] sync [<external>]
  todo [(-c <cfg>) -v] show <id>
  todo [(-c <cfg>) -v] do <id>
  todo [(-c <cfg>) -v] wait <id>
  todo [(-c <cfg>) -v] done <id>
    
Options:
  -a          include all tasks [default false]
  -v          be verbose (debug) [default false]
  -c <cfg>    config [default ~/.todo.conf]
  -d          dry run, only print what would be updated [default false]
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
	Repo     string
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

	config, err := readConfig(opts["-c"])
	if err != nil {
		mainLog.Error(err.Error())
		mainLog.Debug("%+v", err)
		return
	}

	repo := repo(config.Repo)
	if repo == nil {
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
	sort.Strings(keys)
	for i, key := range keys {
		if i != 0 {
			res.WriteString(" ")
		}
		res.WriteString(fmt.Sprintf("%s=%v", key, opts[key]))
	}
	return res.String()
}

func repo(path string) todo.RepoBegin {
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

	f, e := os.Open(p)
	if e != nil {
		return c, errors.Wrap(e, "Could not open config file")
	}
	defer f.Close()

	return readConfigToml(f)
}

func readConfigToml(r io.Reader) (config, error) {
	var raw map[string]interface{}
	var c config
	if _, err := toml.DecodeReader(r, &raw); err != nil {
		return config{}, errors.Wrap(err, "Could not read config")
	}

	for key, value := range raw {
		if key == "repo" {
			if uri, ok := value.(string); ok {
				c.Repo = uri
			} else {
				return c, errors.New("Invalid config, 'repo' should be uri")
			}
		} else {
			if values, ok := value.(map[string]interface{}); ok {
				e := ext.ExternalConfig{
					ID:    key,
					Extra: map[string]string{},
				}
				for valueKey, valueValue := range values {
					if valueKey == "uri" {
						e.URI = valueValue.(string)
					} else if valueKey == "type" {
						e.Type = valueValue.(string)
					} else {
						e.Extra[valueKey] = valueValue.(string)
					}
				}
				if e.ID == "" || e.Type == "" || e.URI == "" {
					return c, errors.New("Invalid config, id, type and uri is required for '" + key + "' ")
				}
				c.External = append(c.External, e)
			} else {
				return c, errors.New("Invalid config, '" + key + "' (external config) should be table")
			}
		}
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
