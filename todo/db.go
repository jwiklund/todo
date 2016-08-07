package todo

import (
	"database/sql"
	"encoding/json"
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

	_, err = db.Exec(`create table if not exists todo(id text primary key, state text, message text, attr text)`)
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
	rows, err := d.db.Query("select rowid, state, message, attr from todo where state != 'done'")
	if err != nil {
		return nil, errors.New(err.Error())
	}
	defer rows.Close()
	var message string
	var state string
	var rowid int64
	var attrB []byte
	var tasks []Task
	for rows.Next() {
		err = rows.Scan(&rowid, &state, &message, &attrB)
		if err != nil {
			return nil, errors.Wrap(err, "Could not scan task")
		}
		log.Debugf("list id=%v,state=%s,message=%s,attr=%s", rowid, state, message, string(attrB))
		attrM, err := decodeAttr(attrB)
		if err != nil {
			return nil, errors.Wrap(err, "Could not decode attribute")
		}
		tasks = append(tasks, Task{
			ID:      strconv.FormatInt(rowid, 10),
			State:   StateFrom(state),
			Message: message,
			Attr:    attrM,
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
	rows, err := d.db.Query("select rowid, state, message, attr from todo where rowid = ?", id)
	if err != nil {
		return Task{}, errors.New(err.Error())
	}
	defer rows.Close()
	if rows.Next() {
		var rowid int64
		var state string
		var message string
		var attrB []byte
		err = rows.Scan(&rowid, &state, &message, &attrB)
		log.Debugf("get id=%v,state=%s,message=%s,attr=%s", rowid, state, message, string(attrB))
		if err != nil {
			return Task{}, errors.New(err.Error())
		}
		attr, err := decodeAttr(attrB)
		if err != nil {
			return Task{}, errors.Wrap(err, "Could not decode attributes")
		}
		return Task{
			ID:      strconv.FormatInt(rowid, 10),
			State:   StateFrom(state),
			Message: message,
			Attr:    attr,
		}, nil
	}
	return Task{}, errors.New("Task not found")
}

func (d *dbRepo) Update(t Task) error {
	attr, err := encodeAttr(t.Attr)
	if err != nil {
		return errors.Wrap(err, "Could not encode attributes")
	}
	log.Debugf("update id=%v,state=%s,message=%s,attr=%s", t.ID, t.State.String(), t.Message, string(attr))
	r, err := d.db.Exec("update todo set state = ?, message = ?, attr = ? where rowid = ?",
		t.State.String(), t.Message, attr, t.ID)
	if err != nil {
		return errors.Wrap(err, "Could not update task")
	}
	if rows, _ := r.RowsAffected(); rows != 1 {
		return errors.New("Update failed, no rows affected")
	}
	return nil
}

func encodeAttr(a map[string]string) ([]byte, error) {
	if a == nil || len(a) == 0 {
		return nil, nil
	}
	attr := struct {
		Attributes map[string]string
	}{a}
	return json.Marshal(&attr)
}

func decodeAttr(bs []byte) (map[string]string, error) {
	if bs == nil || len(bs) == 0 {
		return make(map[string]string), nil
	}
	attr := struct {
		Attributes map[string]string
	}{}
	err := json.Unmarshal(bs, &attr)
	if err != nil {
		return nil, err
	}
	if attr.Attributes == nil {
		return make(map[string]string), nil
	}
	return attr.Attributes, err
}
