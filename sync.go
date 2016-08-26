package main

import (
	_ "github.com/jwiklund/todo/ext/text"
	"github.com/jwiklund/todo/view"
)

// todo [-v][-r <repo>] sync [<external>]
func syncCmd(t view.Todo, opts map[string]interface{}) {
	external := ""
	if s := opts["<external>"]; s != nil {
		external = s.(string)
	}
	dryRun := false
	if b := opts["-d"]; b != nil {
		dryRun = b.(bool)
	}
	if external != "" {
		err := t.Sync(external, dryRun)
		if err != nil {
			mainLog.Errorf("Failed to sync external %s %v", external, err)
			mainLog.Debugf("%+v", err)
			return
		}
	} else {
		err := t.SyncAll(dryRun)
		if err != nil {
			mainLog.Errorf("Failed to sync external %v", err)
			mainLog.Debugf("%+v", err)
			return
		}
	}

	list(t, false, "")
}
