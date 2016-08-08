package text

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/jwiklund/todo/ext"
	"github.com/jwiklund/todo/util"
)

var log = logrus.WithField("comp", "ext.text")

func init() {
	ext.Register("text", func(cfg ext.ExternalConfig) (ext.External, error) {
		return New(util.Expand(cfg.URI))
	})
}

// New create a new file mirror
func New(path string) (ext.External, error) {
	_, err := os.Stat(path)
	if err != nil {
		log.Debugf("mirror does not exist %v", err)
	}
	return nil, fmt.Errorf("not implemented")
}
