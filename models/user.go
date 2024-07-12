package models

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	cUserTable   = "user"
	cUserCId     = "id"
	cUserCName   = "name"
	cUserCUserId = "user_id"
	cUserCUuid   = "uuid"
)

type User struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	UserId string `db:"user_id" json:"userId"`
	Uuid   string `db:"uuid" json:"uuid"`
}

func (u *User) Get(id int64) error {
	return WrapGet(u, id,
		func(wgo *wrapSqlOptions) { wgo.tableName = cUserTable },
	)
}
func (u *User) GetByUserId(userid string) error {
	return WrapGet(u, userid,
		func(wgo *wrapSqlOptions) { wgo.tableName = cUserTable },
		func(wgo *wrapSqlOptions) { wgo.idName = cUserCUserId },
	)
}

func (u *User) Add() error {
	u.Uuid = uuid.NewString()
	sb := from(cUserTable).
		Insert().
		Cols(cUserCName, cUserCUserId, cUserCUuid).
		Vals(goqu.Vals{u.Name, u.UserId, u.Uuid})
	return AddQ(u, sb, &u.Id)
}

func (u *User) Remove() error {
	return WrapRemove(u, u.Id,
		func(wgo *wrapSqlOptions) { wgo.tableName = cUserTable },
	)
}

type Users []User

func (u *Users) All() error {
	sb := from(cUserTable).
		Select(cAll)
	return SelectQ(u, sb)
}
func (u *Users) AllT(tx *sqlx.Tx) error {
	sb := from(cUserTable).
		Select(cAll)
	return SelectQ(u, sb, tx)
}
