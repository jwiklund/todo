package main

import (
	"bytes"
	"fmt"

	log "github.com/Sirupsen/logrus"
	opt "github.com/docopt/docopt-go"
	_ "github.com/pkg/errors"
)

/*

 */
var usage = `Todo list.

Usage:
  todo -h
  todo [-av][-r <repo>] list
  todo [-av][-r <repo>] add <message...>
  todo [-av][-r <repo>] update <id> <state>
    
Options:
  -a          include all tasks [default false]
  -v          be verbose (debug) [default false]
  -r <repo>   custom repo [default sqlite:~/todo.sqlite]
`

func main() {
	opts, err := opt.Parse(usage, nil, true, "1.0", false)
	if err != nil {
		log.Fatal(err)
		return
	}
	if opts["-v"].(bool) {
		log.SetLevel(log.DebugLevel)
	}
	mainLog := log.WithField("comp", "main")
	mainLog.Debug("Args ", sortOpts(opts))
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
