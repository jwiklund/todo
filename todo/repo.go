package todo

import (
	"database/sql"
	"os"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	// load sqlite3 driver needed for sqlite store
	_ "github.com/mattn/go-sqlite3"
)

var (
	log = logrus.WithField("comp", "todo")
)

// Task a todo task
type Task struct {
	ID      string
	Message string
}

func (t Task) String() string {
	return "(" + t.ID + ")\t" + t.Message
}

// IsCurrent return true if task is not waiting or archived
func (t Task) IsCurrent() bool {
	return true
}

// Repo a todo repository
type Repo interface {
	Close() error
	List() ([]Task, error)
	Add(string) (Task, error)
}

// RepoFromPath Return a repository stored in path
func RepoFromPath(path string) (Repo, error) {
	t, p := splitPath(path)
	log.Debugf("Open %s at %s", t, p)
	if t != "sqlite" {
		return nil, errors.Errorf("Invalid repo type %s", t)
	}

	db, err := sql.Open("sqlite3", p)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	_, err = db.Exec(`create table if not exists todo(id text primary key, message text)`)
	if err != nil {
		db.Close()
		return nil, errors.New(err.Error())
	}

	return &dbRepo{db: db}, nil
}

func splitPath(path string) (string, string) {
	t := "sqlite"
	p := path
	indx := strings.Index(path, "://")
	if indx != -1 {
		t = path[0:indx]
		p = path[indx+3:]
	}
	if strings.HasPrefix(p, "~/") {
		home := os.ExpandEnv("$HOME")
		p = home + p[1:]
	}

	return t, p
}

type dbRepo struct {
	db *sql.DB
}

func (d *dbRepo) Close() error {
	if d.db != nil {
		db := d.db
		d.db = nil
		return db.Close()
	}
	return nil
}

func (d *dbRepo) List() ([]Task, error) {
	rows, err := d.db.Query("select rowid, message from todo")
	if err != nil {
		return nil, errors.New(err.Error())
	}
	defer rows.Close()
	var message string
	var rowid int64
	var tasks []Task
	for rows.Next() {
		err = rows.Scan(&rowid, &message)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		tasks = append(tasks, Task{
			ID:      strconv.FormatInt(rowid, 10),
			Message: message,
		})
	}
	return tasks, nil
}

func (d *dbRepo) Add(message string) (Task, error) {
	r, err := d.db.Exec("insert into todo(message) values (?)", message)
	if err != nil {
		return Task{}, errors.New(err.Error())
	}
	id, err := r.LastInsertId()
	if err != nil {
		return Task{}, errors.New(err.Error())
	}
	return Task{
		ID:      strconv.FormatInt(id, 10),
		Message: message,
	}, nil
}
