package config

import (
	"fmt"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
)

func from(table string) *sqlbuilder.CreateTableBuilder {
	return sqlbuilder.CreateTable(table)
}

var dialects = map[string]map[string]string{
	"ramsql": map[string]string{
		"primary": " INT PRIMARY KEY AUTOINCREMENT ",
	},
	"sqlite": map[string]string{
		"primary": " INTEGER PRIMARY KEY AUTOINCREMENT ",
	},
}

func getPrimaryConstant(driver string) string {
	dialects := dialects[driver]
	return dialects["primary"]
}

func CreateUserSql(driver string) string {
	return from("user").IfNotExists().
		Define("id" + getPrimaryConstant(driver)).
		Define("name VARCHAR NOT NULL").
		Define("user_id VARCHAR NOT NULL UNIQUE").
		Define("uuid VARCHAR NOT NULL").
		String()
}
func commonContentSql(driver string) *sqlbuilder.CreateTableBuilder {
	return from("content_request").IfNotExists().
		Define("id" + getPrimaryConstant(driver)).
		Define("user_id INT NOT NULL").
		Define("content VARCHAR NOT NULL").
		Define("uuid VARCHAR NOT NULL").
		Define("create_at DATETIME NOT NULL")
}
func CreateContentSql(driver string) string {
	return commonContentSql(driver).String()
}
func CreateContentRquestSql(driver string) string {
	return commonContentSql(driver).
		Define("uploaded BOOLEAN NOT NULL").
		String()
}

func CreateImageSql(driver string) string {
	return from("content_request").IfNotExists().
		Define("id" + getPrimaryConstant(driver)).
		Define("content_id INT NOT NULL").
		Define("uuid VARCHAR NOT NULL").
		String()
}

func MustCreateDb(driver string) {
	mustExeSqls(
		CreateUserSql(driver),
		CreateContentRquestSql(driver),
		CreateContentSql(driver),
		CreateImageSql(driver),
	)
}
func ShowTables() {
	data, err := GetDb().Queryx(`
		SELECT * 
			FROM sqlite_master 
			WHERE type='table'
	`)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	fmt.Printf("data: %v\n", rowsToMaps(data))
}

func mustExeSqls(sqls ...string) {
	db := GetDb().MustBegin()
	defer MustRecoverTx(db)
	for _, sql := range sqls {
		Logger.Sugar().Debugf("sql: %s, args: %v", sql, nil)
		db.MustExec(sql)
	}
}

func rowsToMaps(rs *sqlx.Rows) any {
	var res []map[string]interface{}
	for rs.Next() {
		d := make(map[string]interface{})
		rs.MapScan(d)
		res = append(res, d)
	}
	return res
}