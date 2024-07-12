package config

import (
	"context"
	"database/sql"
	"os"
	"sync"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	_ "github.com/proullon/ramsql/driver"
	_ "modernc.org/sqlite"
)

const DbConfigGoquDiarect = "sqlite3"
const DbConfig = "./tmp/sqlite.db"
const DbConfigDriver = "sqlite"

// var Db *sqlx.DB
var Db = sync.OnceValue(func() *sqlx.DB {
	db, err := sqlx.Open(DbConfigDriver, DbConfig)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	Logger.Debug("start db")
	return db
})

func GetDb() *sqlx.DB {
	return Db()
}
func DeferDb() {
	GetDb().Close()
}
func BeginTxx() (*sqlx.Tx, error) {
	return GetDb().BeginTxx(context.Background(), &sql.TxOptions{})
}

func RecoverTx(db *sqlx.Tx) {
	if err := recover(); err != nil {
		db.Rollback()
	} else {
		db.Commit()
	}
}
func MustRecoverTx(db *sqlx.Tx) {
	if err := recover(); err != nil {
		db.Rollback()
		panic(err)
	} else {
		db.Commit()
	}
}

var goqud = goqu.Dialect(DbConfigGoquDiarect)

func Goqu() goqu.DialectWrapper {
	return goqud
}

func SetTestDb() {
	// driver, url := DbConfigDriver, "../tmp/test.db"
	driver, url := "ramsql", "test"
	if driver == "sqlite" {
		f, _ := os.Create(url)
		f.Close()
	}
	Db = sync.OnceValue(func() *sqlx.DB {
		db, err := sqlx.Open(driver, url)
		// db, err := sqlx.Open(dirver, url)
		if err != nil {
			panic(err)
		}
		err = db.Ping()
		if err != nil {
			panic(err)
		}
		Logger.Debug("start test db")
		return db
	})
	MustCreateDb(driver)
}
