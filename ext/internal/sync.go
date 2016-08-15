package internal

import (
	"github.com/Sirupsen/logrus"
	"github.com/deckarep/golang-set"
	"github.com/jwiklund/todo/todo"
)

type indexedTasks struct {
	tasks []todo.Task
	index map[string]int
}

// SyncHelper common sync implememntation
func SyncHelper(r todo.RepoBegin, extID string, dryRun bool, externalCurrent, localCurrent []todo.Task) error {
	external := index(extID, externalCurrent)
	local := index(extID, localCurrent)

	externalIds := external.IDSet()
	localIds := local.IDSet()

	added := externalIds.Difference(localIds)
	missing := localIds.Difference(externalIds)
	updated := updated(external, local)

	if dryRun {
		syncLog := logrus.WithField("comp", "ext.sync")
		syncLog.Info("DRY RUN ", extID, " ADD ", added)
		syncLog.Info("DRY RUN ", extID, " REM ", missing)
		syncLog.Info("DRY RUN ", extID, " UPD ", updated)
		return nil
	}

	commitRepo, err := r.Begin()
	if err != nil {
		return err
	}
	if err := syncAdd(commitRepo, extID, external, added); err != nil {
		commitRepo.Close()
		return err
	}
	if err := syncRemove(commitRepo, local, missing); err != nil {
		commitRepo.Close()
		return err
	}
	if err := syncUpdate(commitRepo, external, local, updated); err != nil {
		commitRepo.Close()
		return err
	}
	return commitRepo.Commit()
}

func index(extID string, tasks []todo.Task) *indexedTasks {
	index := map[string]int{}
	for i, task := range tasks {
		if a, ok := task.Attr[extID+".id"]; ok {
			index[a] = i
		}
	}
	return &indexedTasks{tasks, index}
}

func (i *indexedTasks) IDSet() mapset.Set {
	r := mapset.NewThreadUnsafeSet()
	for key := range i.index {
		r.Add(key)
	}
	return r
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

func updated(external, local *indexedTasks) mapset.Set {
	updated := mapset.NewThreadUnsafeSet()

	for externalKey, externalIndex := range external.index {
		if localIndex, ok := local.index[externalKey]; ok {
			if external.tasks[externalIndex].Message != local.tasks[localIndex].Message {
				updated.Add(externalKey)
			}
		}
	}

	return updated
}

func syncAdd(r todo.Repo, extID string, external *indexedTasks, added mapset.Set) error {
	for _, add := range added.ToSlice() {
		index := external.index[add.(string)]
		addedMessage := external.tasks[index].Message
		_, err := r.AddWithAttr(addedMessage, map[string]string{
			"external":    extID,
			extID + ".id": add.(string),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func syncRemove(r todo.Repo, local *indexedTasks, missing mapset.Set) error {
	for _, rem := range missing.ToSlice() {
		index := local.index[rem.(string)]
		t := local.tasks[index]
		t.State = todo.StateFrom("done")
		err := r.Update(t)
		if err != nil {
			return err
		}
	}
	return nil
}

func syncUpdate(r todo.Repo, external, local *indexedTasks, updated mapset.Set) error {
	for _, up := range updated.ToSlice() {
		externalIndex := external.index[up.(string)]
		localIndex := local.index[up.(string)]
		e := external.tasks[externalIndex]
		l := local.tasks[localIndex]
		l.Message = e.Message
		err := r.Update(l)
		if err != nil {
			return err
		}
	}
	return nil
}
