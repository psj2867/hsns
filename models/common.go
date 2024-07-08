package models

import (
	"database/sql"

	"github.com/psj2867/hsns/config"
)

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
func ExecQ[T any](t *T, sb goquSql, dbs ...sqlxDb) (sql.Result, error) {
	sql, args, err := sb.ToSQL()
	if err != nil {
		return nil, err
	}
	config.Logger.Sugar().Debugf("sql: %s, args: %v", sql, args)
	return getDb(dbs).Exec(sql, args...)
}

func Refresh(t *int64, r sql.Result) error {
	l, err := r.LastInsertId()
	if err != nil {
		return err
	}
	(*t) = l
	return nil
}
