// +build !windows

package util

import "os"

func home() string {
	return os.ExpandEnv("$HOME") + "/"
}
