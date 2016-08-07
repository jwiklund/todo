package todo

import (
	"database/sql"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	// load sqlite3 driver needed for sqlite store
	_ "github.com/mattn/go-sqlite3"
)

func newSQL(t, p string) (Repo, error) {
	db, err := sql.Open("sqlite3", p)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	_, err = db.Exec(`create table if not exists todo(id text primary key, state text, message text)`)
	if err != nil {
		db.Close()
		return nil, errors.New(err.Error())
	}
	return &dbRepo{db}, nil
}

func splitPath(path string) (string, string) {
	t := "sqlite"
	p := path
	indx := strings.Index(path, "://")
	if indx != -1 {
		t = path[0:indx]
		p = path[indx+3:]
	}
	log.Debug(os.Environ())
	if strings.HasPrefix(p, "~/") {
		p = home() + p[2:]
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
