package todo

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/jwiklund/todo/util"
	"github.com/pkg/errors"
	// load sqlite3 driver needed for sqlite store
	_ "github.com/mattn/go-sqlite3"
)

func newSQL(t, p string) (RepoBegin, error) {
	db, err := sql.Open("sqlite3", p)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	_, err = db.Exec(`create table if not exists todo(
		state text, message text,
		repo text,
		ext_id text,
		attr text)`)
	if err != nil {
		db.Close()
		return nil, errors.Wrap(err, "Could not create task table")
	}
	_, err = db.Exec(`create unique index if not exists external_idx on todo(repo, ext_id)`)
	if err != nil {
		db.Close()
		return nil, errors.Wrap(err, "Could not create external index")
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

	return t, util.Expand(p)
}

type dbOrTx interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
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

func (d *dbRepo) Begin() (RepoCommit, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "Could not start transaction")
	}
	return &txRepo{tx}, nil
}

func (d *dbRepo) List() ([]Task, error) {
	return list(d.db)
}

func (d *dbRepo) Add(message string) (Task, error) {
	return add(d.db, message, nil)
}

func (d *dbRepo) AddWithAttr(message string, attr map[string]string) (Task, error) {
	return add(d.db, message, attr)
}

func (d *dbRepo) Update(task Task) error {
	return update(d.db, task)
}

func (d *dbRepo) Get(id string) (Task, error) {
	return get(d.db, id)
}

func (d *dbRepo) GetByExternal(repo, extID string) (Task, error) {
	return getByExternal(d.db, repo, extID)
}

type txRepo struct {
	tx *sql.Tx
}

func (t *txRepo) Begin() (RepoCommit, error) {
	return t, nil
}

func (t *txRepo) Commit() error {
	return t.tx.Commit()
}

func (t *txRepo) Close() error {
	return t.tx.Rollback()
}

func (t *txRepo) List() ([]Task, error) {
	return list(t.tx)
}

func (t *txRepo) Add(message string) (Task, error) {
	return add(t.tx, message, nil)
}

func (t *txRepo) AddWithAttr(message string, attr map[string]string) (Task, error) {
	return add(t.tx, message, attr)
}

func (t *txRepo) Update(task Task) error {
	return update(t.tx, task)
}

func (t *txRepo) Get(id string) (Task, error) {
	return get(t.tx, id)
}

func (t *txRepo) GetByExternal(repo, extID string) (Task, error) {
	return getByExternal(t.tx, repo, extID)
}

func list(db dbOrTx) ([]Task, error) {
	rows, err := db.Query("select rowid, state, message, attr from todo where state != 'done'")
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
		todoLog.Debugf("list id=%v,state=%s,message=%s,attr=%s", rowid, state, message, string(attrB))
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

func add(db dbOrTx, message string, attr map[string]string) (Task, error) {
	attrB, err := encodeAttr(attr)
	if err != nil {
		return Task{}, errors.Wrap(err, "could not encode attr")
	}
	repo, extID := getExternal(attr)
	todoLog.Debugf("add state=%s,message=%s,repo=%s,ext_id=%s,attr=%s",
		"todo", message, nullable(repo), nullable(extID), string(attrB))
	r, err := db.Exec("insert into todo(state, message, repo, ext_id, attr) values (?, ?, ?, ?, ?)",
		"todo", message, repo, extID, attrB)
	if err != nil {
		return Task{}, errors.Wrap(err, "could not write task")
	}
	id, err := r.LastInsertId()
	if err != nil {
		return Task{}, errors.Wrap(err, "could not get id")
	}
	return Task{
		ID:      strconv.FormatInt(id, 10),
		State:   "todo",
		Message: message,
		Attr:    attr,
	}, nil
}

func getByRows(rows *sql.Rows) (Task, error) {
	if rows.Next() {
		var rowid int64
		var state string
		var message string
		var attrB []byte
		err := rows.Scan(&rowid, &state, &message, &attrB)
		if err != nil {
			return Task{}, errors.Wrap(err, "Could not scan row")
		}
		todoLog.Debugf("get id=%v,state=%s,message=%s,attr=%s", rowid, state, message, string(attrB))
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
	return Task{}, ErrorNotFound
}

func get(db dbOrTx, id string) (Task, error) {
	query := `select rowid, state, message, attr
	            from todo
			   where rowid = ?`
	rows, err := db.Query(query, id)
	if err != nil {
		return Task{}, errors.New(err.Error())
	}
	defer rows.Close()
	return getByRows(rows)
}

func getByExternal(db dbOrTx, repo, extID string) (Task, error) {
	query := `select rowid, state, message, attr 
	            from todo 
	           where repo = ? and ext_id = ?`
	rows, err := db.Query(query, repo, extID)
	if err != nil {
		return Task{}, errors.Wrap(err, "Could not query")
	}
	defer rows.Close()
	return getByRows(rows)
}

func update(db dbOrTx, t Task) error {
	attr, err := encodeAttr(t.Attr)
	if err != nil {
		return errors.Wrap(err, "Could not encode attributes")
	}
	repo, extID := getExternal(t.Attr)
	todoLog.Debugf("update id=%v,state=%s,message=%s,repo=%s,ext_id=%s,attr=%s",
		t.ID, t.State.String(), t.Message, nullable(repo), nullable(extID), string(attr))
	r, err := db.Exec("update todo set state = ?, message = ?, repo = ?, ext_id = ?, attr = ? where rowid = ?",
		t.State.String(), t.Message, repo, extID, attr, t.ID)
	if err != nil {
		return errors.Wrap(err, "Could not update task")
	}
	if rows, _ := r.RowsAffected(); rows != 1 {
		return errors.New("Update failed, no rows affected")
	}
	return nil
}
func getExternal(a map[string]string) (*string, *string) {
	repo, ok := a["external"]
	if !ok {
		return nil, nil
	}
	extID, ok := a[repo+".id"]
	if !ok {
		return nil, nil
	}
	return &repo, &extID
}

func nullable(s *string) string {
	if s == nil {
		return "<nil>"
	}
	return *s
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
