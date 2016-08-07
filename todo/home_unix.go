// +build !windows

package todo

import "os"

func home() string {
	return os.ExpandEnv("$HOME") + "/"
}
