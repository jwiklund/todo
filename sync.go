package main

import (
	"github.com/jwiklund/todo/ext"
	_ "github.com/jwiklund/todo/ext/text"
)

// todo [-v][-r <repo>] sync [<external>]
func syncCmd(r ext.Repo, opts map[string]interface{}) {
	external := ""
	if s := opts["<external>"]; s != nil {
		external = s.(string)
	}
	if external != "" {
		err := r.Sync(external)
		if err != nil {
			mainLog.Errorf("Failed to sync external %s %v", external, err)
			mainLog.Debugf("%+v", err)
			return
		}
	} else {
		err := r.SyncAll()
		if err != nil {
			mainLog.Errorf("Failed to sync external %v", err)
			mainLog.Debugf("%+v", err)
			return
		}
	}

	list(r, false, "")
}
