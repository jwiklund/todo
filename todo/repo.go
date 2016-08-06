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

// State valid state, one of "todo", "doing", "waiting", done"
type State string

func (s State) String() string {
    return string(s)
}

// StateValid check if state is a valid state
func StateValid(state string) bool {
    return state == StateFrom(state).String()
}
// StateFrom returns a valid state from a given string (or todo)
func StateFrom(state string) State {
    switch state {
        case "todo":
            return State(state)
        case "doing":
            return State(state)
        case "waiting":
            return State(state)
        case "done":
            return State(state)
    }
    return "todo"
}

// Task a todo task
type Task struct {
	ID      string
	State   State
	Message string
}

func (t Task) String() string {
	return "(" + t.ID + ")\t" + t.Message
}

// IsCurrent return true if task is not waiting or archived
func (t Task) IsCurrent() bool {
	return t.State == "todo" || t.State == "doing"
}

// Repo a todo repository
type Repo interface {
	Close() error
	List() ([]Task, error)
	Add(string) (Task, error)
	Get(string) (Task, error)
	Update(Task) error
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

	_, err = db.Exec(`create table if not exists todo(id text primary key, state text, message text)`)
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
	rows, err := d.db.Query("select rowid, state, message from todo where state != 'done'")
	if err != nil {
		return nil, errors.New(err.Error())
	}
	defer rows.Close()
	var message string
	var state string
	var rowid int64
	var tasks []Task
	for rows.Next() {
		err = rows.Scan(&rowid, &state, &message)
		if err != nil {
			return nil, errors.New(err.Error())
		}
        log.Debugf("raw id=%v,state=%s,message=%s", rowid, state, message)
		tasks = append(tasks, Task{
			ID:      strconv.FormatInt(rowid, 10),
			State:   StateFrom(state),
			Message: message,
		})
	}
	return tasks, nil
}

func (d *dbRepo) Add(message string) (Task, error) {
	r, err := d.db.Exec("insert into todo(state, message) values (?, ?)", "todo", message)
	if err != nil {
		return Task{}, errors.New(err.Error())
	}
	id, err := r.LastInsertId()
	if err != nil {
		return Task{}, errors.New(err.Error())
	}
	return Task{
		ID:      strconv.FormatInt(id, 10),
		State:   "todo",
		Message: message,
	}, nil
}

func (d *dbRepo) Get(id string) (Task, error) {
	rows, err := d.db.Query("select rowid, state, message from todo where rowid = ?", id)
	if err != nil {
		return Task{}, errors.New(err.Error())
	}
	defer rows.Close()
	if rows.Next() {
		var id int64
		var state string
		var message string
		err = rows.Scan(&id, &state, &message)
		if err != nil {
			return Task{}, errors.New(err.Error())
		}
		return Task{
			ID:      strconv.FormatInt(id, 10),
			State:   StateFrom(state),
			Message: message,
		}, nil
	}
	return Task{}, errors.New("Task not found")
}

func (d *dbRepo) Update(t Task) error {
	r, err := d.db.Exec("update todo set state = ?, message = ? where rowid = ?", t.State.String(), t.Message, t.ID)
	if err != nil {
		return errors.New(err.Error())
	}
	if rows, _ := r.RowsAffected(); rows != 1 {
		return errors.New("Update failed, no rows affected")
	}
	return nil
}
