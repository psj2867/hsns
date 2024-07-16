package config

import (
	"context"
	"database/sql"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	_ "github.com/proullon/ramsql/driver"
	"github.com/stoewer/go-strcase"
	_ "modernc.org/sqlite"
)

const DbConfigGoquDiarect = "sqlite3"
const DbConfig = "./tmp/sqlite.db"
const DbConfigDriver = "sqlite"

func goquColumnRename(org string) string {
	return strcase.SnakeCase(org)
}
func setGoqu() {
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		panic(err)
	}
	goqu.SetDefaultPrepared(true)
	goqu.SetColumnRenameFunction(goquColumnRename)
	goqu.SetTimeLocation(loc)
}

// var Db *sqlx.DB
var Db = sync.OnceValue(func() *sqlx.DB {
	setGoqu()
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

func SetTestDb(driver string, url string) {
	if strings.HasPrefix(driver, "sqlite") {
		f, _ := os.Create(url)
		f.Close()
	}
	Db = sync.OnceValue(func() *sqlx.DB {
		setGoqu()
		db, err := sqlx.Open(driver, url)
		if err != nil {
			panic(err)
		}
		err = db.Ping()
		if err != nil {
			panic(err)
		}
		Logger.Sugar().Debugf("start test db driver: %s, url: %s", driver, url)
		return db
	})
	MustCreateDb(driver)
}
