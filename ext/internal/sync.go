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

	rev, err := revived(r, extID, external, added)
	if err != nil {
		return err
	}
	for _, r := range rev {
		added.Remove(r.Attr[extID+".id"])
	}

	if dryRun {
		syncLog := logrus.WithField("comp", "ext.sync")
		syncLog.Infof("DRY RUN EXT %+v", externalCurrent)
		syncLog.Infof("DRY RUN LOC %+v", localCurrent)
		syncLog.Info("DRY RUN ", extID, " ADD ", added)
		syncLog.Info("DRY RUN ", extID, " REM ", missing)
		syncLog.Info("DRY RUN ", extID, " UPD ", updated)
		syncLog.Info("DRY RUN ", extID, " REV ", rev)
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
	if err := syncRevive(commitRepo, rev); err != nil {
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

func (i *indexedTasks) GetByExternal(extID string) todo.Task {
	index, _ := i.index[extID]
	return i.tasks[index]
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
			e := external.tasks[externalIndex]
			l := local.tasks[localIndex]

			if e.Message != e.Message || e.State != l.State {
				updated.Add(externalKey)
			}

		}
	}

	return updated
}

func merge(external, local todo.Task) todo.Task {
	for key, value := range local.Attr {
		if _, ok := external.Attr[key]; !ok {
			external.Attr[key] = value
		}
	}
	external.ID = local.ID
	return external
}

func revived(r todo.Repo, extID string, external *indexedTasks, missing mapset.Set) ([]todo.Task, error) {
	var result []todo.Task

	for _, removed := range missing.ToSlice() {
		r, err := r.GetByExternal(extID, removed.(string))
		if err != nil {
			if err != todo.ErrorNotFound {
				return nil, err
			}
			// not found, ignore
		} else {
			// found, not revived
			result = append(result, merge(external.GetByExternal(removed.(string)), r))
		}
	}

	return result, nil
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
		l.State = e.State
		err := r.Update(l)
		if err != nil {
			return err
		}
	}
	return nil
}

func syncRevive(r todo.Repo, revived []todo.Task) error {
	for _, t := range revived {
		err := r.Update(t)
		if err != nil {
			return err
		}
	}
	return nil
}
