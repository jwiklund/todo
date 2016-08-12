package text

import (
	"strconv"

	"github.com/deckarep/golang-set"
	"github.com/jwiklund/todo/todo"
)

// Sync all available lines with all available items
// TODO find existing done tasks before creating new tasks (needs find by attribute)
func (t *text) Sync(r todo.RepoBegin) error {
	tasks, err := r.List()
	if err != nil {
		return err
	}

	rowToTask := rowToTask(t.id, tasks)
	current := currentSet(rowToTask)
	infile := infileSet(t.source)

	added := infile.Difference(current)
	missing := current.Difference(infile)
	updated := updated(t.id, t.source, rowToTask, tasks)

	commitRepo, err := r.Begin()
	if err != nil {
		return err
	}
	if err := syncAdd(commitRepo, t.id, t.source, added); err != nil {
		commitRepo.Close()
		return err
	}
	if err := syncRemove(commitRepo, rowToTask, tasks, missing); err != nil {
		commitRepo.Close()
		return err
	}
	if err := syncUpdate(commitRepo, t.source, tasks, rowToTask, updated); err != nil {
		commitRepo.Close()
		return err
	}
	return commitRepo.Commit()
}

func rowToTask(id string, tasks []todo.Task) map[int]int {
	rowToTask := map[int]int{}
	for tasknr, task := range tasks {
		if attr, ok := task.Attr[id+".id"]; ok {
			attrint, err := strconv.Atoi(attr)
			if err != nil {
				textLog.Errorf("Invalid external id for task %s", task.ID)
			} else {
				rowToTask[attrint] = tasknr
			}
		}
	}
	return rowToTask
}

func currentSet(mapping map[int]int) mapset.Set {
	s := mapset.NewThreadUnsafeSet()
	for attr := range mapping {
		s.Add(attr)
	}
	return s
}

func infileSet(lines [][]byte) mapset.Set {
	s := mapset.NewThreadUnsafeSet()
	for index := range lines {
		s.Add(index)
	}
	return s
}

func updated(id string, source [][]byte, rowToTask map[int]int, tasks []todo.Task) mapset.Set {
	updated := mapset.NewThreadUnsafeSet()
	for rowIndex, bytes := range source {
		if taskIndex, ok := rowToTask[rowIndex]; ok {
			if string(bytes) != tasks[taskIndex].Message {
				updated.Add(rowIndex)
			}
		}
	}

	return updated
}

func syncAdd(r todo.Repo, id string, source [][]byte, added mapset.Set) error {
	for index, add := range added.ToSlice() {
		addedMessage := string(source[add.(int)])
		_, err := r.AddWithAttr(addedMessage, map[string]string{
			"external": id,
			id + ".id": strconv.Itoa(index),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func syncRemove(r todo.Repo, rowToTask map[int]int, tasks []todo.Task, missing mapset.Set) error {
	for _, rem := range missing.ToSlice() {
		if tindex, ok := rowToTask[rem.(int)]; ok {
			t := tasks[tindex]
			t.State = todo.StateFrom("done")
			err := r.Update(t)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func syncUpdate(r todo.Repo, source [][]byte, tasks []todo.Task, rowToTask map[int]int, updated mapset.Set) error {
	for _, up := range updated.ToSlice() {
		if tindex, ok := rowToTask[up.(int)]; ok {
			t := tasks[tindex]
			t.Message = string(source[up.(int)])
			err := r.Update(t)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
