package ext

import (
	"github.com/jwiklund/todo/todo"
)

// MapExternal return a map from external id -> task id
func MapExternal(external string, tasks []todo.Task) map[string]string {
	res := map[string]string{}
	for _, t := range tasks {
		if attr, ok := t.Attr[external+".id"]; ok {
			res[attr] = t.ID
		}
	}
	return res
}
