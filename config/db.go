package config

import (
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

var Db *sqlx.DB

func init() {
	var err error
	Db, err = sqlx.Open("sqlite", "./tmp/sqlite.db")
	Db.Ping()
	if err != nil {
		panic(err)
	}
	sqlbuilder.DefaultFlavor = sqlbuilder.SQLite

}

func GetDb() *sqlx.DB {
	return Db
}
