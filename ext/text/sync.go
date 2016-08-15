package text

import (
	"strconv"

	"github.com/jwiklund/todo/ext/internal"
	"github.com/jwiklund/todo/todo"
)

// Sync all available lines with all available items
// TODO find existing done tasks before creating new tasks (needs find by attribute)
func (t *text) Sync(r todo.RepoBegin) error {
	localTasks, err := r.List()
	if err != nil {
		return err
	}

	externalTasks := tasksFor(t.id, t.source)

	return internal.SyncHelper(r, t.id, externalTasks, localTasks)

}

func tasksFor(id string, source [][]byte) []todo.Task {
	var res []todo.Task
	for i, message := range source {
		res = append(res, todo.Task{
			Message: string(message),
			Attr: map[string]string{
				"external": id,
				id + ".id": strconv.Itoa(i),
			},
		})
	}
	return res
}
