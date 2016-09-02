package main

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/jwiklund/todo/todo"
)

// Prio render prio
func Prio(prio int) string {
	if prio <= 10 {
		return "high"
	}
	if prio <= 100 {
		return "norm"
	}
	if prio <= 500 {
		return "low"
	}
	if prio <= 1000 {
		return "none"
	}
	return "dont"
}

func renderList(ts []todo.Task, out io.Writer) {
	w := tabwriter.NewWriter(os.Stdout, 6, 8, 1, ' ', 0)
	for _, task := range ts {
		fmt.Fprintf(w, "(%s)\t%s\t%s\t%s\n", task.ID, Prio(task.Prio()), task.State.String(), task.Message)
	}
	w.Flush()
}

func renderOne(task todo.Task, out io.Writer) {
	w := tabwriter.NewWriter(os.Stdout, 6, 8, 2, ' ', 0)
	fmt.Fprintf(w, "(%s)\t%s\t%s\n", task.ID, Prio(task.Prio()), task.State.String(), task.Message)
	for key, value := range task.Attr {
		fmt.Fprintf(w, "\t%s\t%s\n", key, value)
	}
	w.Flush()
}
func render(task todo.Task) string {
	return task.String()
}
