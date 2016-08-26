package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/jwiklund/todo/view"
)

//  todo [-v][-r <repo>] show <id>
func showCmd(t view.Todo, opts map[string]interface{}) {
	task, err := t.Get(opts["<id>"].(string))
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
