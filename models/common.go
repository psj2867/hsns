package models

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"github.com/psj2867/hsns/config"
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
		return dbs[0]
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
func AddQ(sb goquSql, id *int64, dbs ...sqlxDb) error {
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

func defAndSet(options []wrapSqlOptionF) *wrapSqlOptions {
	var w wrapSqlOptions
	w.idName = "id"
	w.tableName = ""
	w.columns = cAll
	for _, v := range options {
		v(&w)
	}
	return &w
}

type wrapSqlOptionF func(*wrapSqlOptions)

func withTx(tx *sqlx.Tx) wrapSqlOptionF { return func(wso *wrapSqlOptions) { wso.tx = tx } }
func withTable(table string) wrapSqlOptionF {
	return func(wso *wrapSqlOptions) { wso.tableName = table }
}

func wrapGet[T any](t *T, id any, options ...wrapSqlOptionF) error {
	get := defAndSet(options)
	sb := from(get.tableName).
		Select(get.columns).
		Where(goqu.C(get.idName).Eq(id))
	if get.tx != nil {
		return GetQ(t, sb, get.tx)
	}
	return GetQ(t, sb)
}

func wrapRemove(id any, options ...wrapSqlOptionF) error {
	get := defAndSet(options)
	sb := from(get.tableName).
		Delete().
		Where(goqu.C(get.idName).Eq(id))
	if get.tx != nil {
		_, err := ExecQ(sb, get.tx)
		return err
	}
	_, err := ExecQ(sb)
	return err
}
