package models

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/psj2867/hsns/config"
)

type goquSql interface {
	ToSQL() (sql string, params []interface{}, err error)
}

type goquStrcutAbsModel[T any] struct{}

func GetQ[T any](t *T, sb *goqu.SelectDataset) error {
	sql, args, err := sb.ToSQL()
	if err != nil {
		return err
	}
	config.Logger.Sugar().Debugf("sql: %s, args: %v", sql, args)
	return config.GetDb().Get(t, sql, args...)
}
func SelectQ[T any](t *T, sb *goqu.SelectDataset) error {
	sql, args, err := sb.ToSQL()
	if err != nil {
		return err
	}
	config.Logger.Sugar().Debugf("sql: %s, args: %v", sql, args)
	return config.GetDb().Select(t, sql, args...)
}
func ExecQ[T any](t *T, sb goquSql) (sql.Result, error) {
	sql, args, err := sb.ToSQL()
	if err != nil {
		return nil, err
	}
	config.Logger.Sugar().Debugf("sql: %s, args: %v", sql, args)
	return config.GetDb().Exec(sql, args...)
}

func Refresh(t *int64, r sql.Result) error {
	l, err := r.LastInsertId()
	if err != nil {
		return err
	}
	(*t) = l
	return nil
}
