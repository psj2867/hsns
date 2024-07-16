package models

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/psj2867/hsns/server/middleware"
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
	return wrapGet(u, id, withTable(cUserTable))
}
func (u *User) GetByUserId(userid string) error {
	return wrapGet(u, userid, withTable(cUserTable),
		func(wgo *wrapSqlOptions) { wgo.idName = cUserCUserId },
	)
}

func (u *User) Add() error {
	u.Uuid = uuid.NewString()
	sb := from(cUserTable).
		Insert().
		Cols(cUserCName, cUserCUserId, cUserCUuid).
		Vals(goqu.Vals{u.Name, u.UserId, u.Uuid})
	return AddQ(sb, &u.Id)
}

func (u *User) Remove() error {
	return wrapRemove(u.Id, withTable(cUserTable))
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

const userContextkey = "psj2867/models/user"

func GetUserInfoInContext(c *gin.Context) *User {
	user, ok := c.Get(userContextkey)
	if !ok {
		setUserInfoInContext(c)
		user, _ = c.Get(userContextkey)
	}
	res := user.(User)
	return &res
}
func setUserInfoInContext(c *gin.Context) {
	userId, ok := middleware.GetAuthUserId(c)
	if !ok {
		return
	}
	user := User{}
	user.GetByUserId(userId)
	c.Set(userContextkey, user)
}
