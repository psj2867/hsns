package models

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"github.com/psj2867/hsns/config"
	_ "modernc.org/sqlite"
)

const (
	cAll = "*"
	cId  = "id"
)

var goqud = config.Goqu()

func from(table string) *goqu.SelectDataset {
	return goqud.From(table)
}

type sqlxDb interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...any) (sql.Result, error)
}

type goquSql interface {
	ToSQL() (sql string, params []interface{}, err error)
}

func getDb(dbs []sqlxDb) sqlxDb {
	if len(dbs) == 0 {
		return config.GetDb()
	} else {
		return dbs[1]
	}

}
func GetQ[T any](t *T, sb goquSql, dbs ...sqlxDb) error {
	sql, args, err := sb.ToSQL()
	if err != nil {
		return err
	}
	config.Logger.Sugar().Debugf("sql: %s, args: %v", sql, args)
	return getDb(dbs).Get(t, sql, args...)
}
func SelectQ[T any](t *T, sb goquSql, dbs ...sqlxDb) error {
	sql, args, err := sb.ToSQL()
	if err != nil {
		return err
	}
	config.Logger.Sugar().Debugf("sql: %s, args: %v", sql, args)
	return getDb(dbs).Select(t, sql, args...)
}
func ExecQ(sb goquSql, dbs ...sqlxDb) (sql.Result, error) {
	sql, args, err := sb.ToSQL()
	if err != nil {
		return nil, err
	}
	config.Logger.Sugar().Debugf("sql: %s, args: %v", sql, args)
	return getDb(dbs).Exec(sql, args...)
}
func AddQ[T any](t *T, sb goquSql, id *int64, dbs ...sqlxDb) error {
	r, err := ExecQ(sb)
	if err != nil {
		return err
	}
	err = Refresh(id, r)
	return err
}

func Refresh(t *int64, r sql.Result) error {
	l, err := r.LastInsertId()
	if err != nil {
		return err
	}
	(*t) = l
	return nil
}

type wrapSqlOptions struct {
	tableName string
	idName    string
	columns   string
	tx        *sqlx.Tx
}

func defAndSet(options []WrapSqlOptionF) *wrapSqlOptions {
	var w wrapSqlOptions
	w.idName = "id"
	w.tableName = ""
	w.columns = cAll
	for _, v := range options {
		v(&w)
	}
	return &w
}

type WrapSqlOptionF func(*wrapSqlOptions)

func WrapGet[T any](t *T, id any, options ...WrapSqlOptionF) error {
	get := defAndSet(options)
	sb := from(get.tableName).
		Select(get.columns).
		Where(goqu.C(get.idName).Eq(id))
	if get.tx != nil {
		return GetQ(t, sb, get.tx)
	}
	return GetQ(t, sb)
}

func WrapRemove[T any](t *T, id int64, options ...WrapSqlOptionF) error {
	get := defAndSet(options)
	sb := from(get.tableName).
		Delete().
		Where(goqu.C(get.idName).Eq(id))
	_, err := ExecQ(sb)
	return err
}
