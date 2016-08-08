package text

import (
	"bytes"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/jwiklund/todo/ext"
	"github.com/jwiklund/todo/todo"
	"github.com/jwiklund/todo/util"
	"github.com/pkg/errors"
)

var log = logrus.WithField("comp", "ext.text")

func init() {
	ext.Register("text", func(cfg ext.ExternalConfig) (ext.External, error) {
		return New(cfg.ID, util.Expand(cfg.URI))
	})
}

// New create a new file mirror
func New(id, path string) (ext.External, error) {
	stat, err := os.Stat(path)
	var source [][]byte
	if err != nil {
		log.Debug("Export file not found", err)
	} else if stat.IsDir() {
		return nil, errors.New("text export is directory")
	} else {
		input, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, errors.Wrap(err, "Could not open exported text")
		}
		source = bytes.Split(input, []byte("\n"))
	}
	return &text{id, path, false, source}, nil
}

type text struct {
	id      string
	path    string
	updated bool
	source  [][]byte
}

func (t *text) Handle(task todo.Task) (todo.Task, error) {
	if a := task.Attr["external"]; a != t.id {
		log.Debug(a, t.id)
		return task, nil
	}
	if a := task.Attr[t.id+".id"]; a != "" {
		ind, err := strconv.Atoi(a)
		if err != nil {
			log.Debugf("Invalid id attribute %s, ignoring %s", a, task.ID)
			return task, nil
		}
		if ind < 0 || ind >= len(t.source) {
			log.Debugf("Invalid id attribute %s (range)", ind)
			return task, nil
		}
		if string(t.source[ind]) != task.Message {
			t.source[ind] = []byte(task.Message)
			t.updated = true
		}
	} else {
		t.source = append(t.source, []byte(task.Message))
		task.Attr[t.id+".id"] = strconv.Itoa(len(t.source) - 1)
		t.updated = true
	}
	return task, nil
}

func (t *text) Close() error {
	if t.updated {
		return ioutil.WriteFile(t.path, bytes.Join(t.source, []byte("\n")), 0660)
	}
	return nil
}
