package util

import "strings"

// Expand ~/path to $HOME/path
func Expand(path string) string {
	if strings.HasPrefix(path, "~/") {
		return home() + path[2:]
	}
	return path
}
