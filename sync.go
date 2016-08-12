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
		if ext, ok := r.Externals()[external]; ok {
			err := ext.Sync(r)
			if err != nil {
				mainLog.Errorf("Failed to sync external %s %v", external, err)
				mainLog.Debugf("%+v", err)
			}
		} else {
			mainLog.Errorf("external %s not configured", external)
		}
	} else {
		for id, ext := range r.Externals() {
			if err := ext.Sync(r); err != nil {
				mainLog.Errorf("Failed to sync external %s %v", id, err)
				mainLog.Debugf("%+v", err)
				return
			}
		}
	}
}
