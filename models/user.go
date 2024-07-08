package models

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/guregu/null/v5"
	"github.com/jmoiron/sqlx"
)

const (
	userTable     = "user"
	userCAll      = "*"
	userCId       = "id"
	userCName     = "name"
	userCFullName = "fullname"
)

type User struct {
	Id       int64
	Name     string
	Fullname null.String
}

func (u *User) Get(id int) error {
	sb := goqu.From(userTable).
		Select(userCAll).
		Where(goqu.Ex{
			userCId: id,
		})
	return GetQ(u, sb)
}

func (u *User) Add() error {
	sb := goqu.From(userTable).
		Insert().
		Cols(userCName, userCFullName).
		Vals(goqu.Vals{u.Name, u.Fullname})
	r, err := ExecQ(u, sb)
	if err != nil {
		return err
	}
	err = Refresh(&u.Id, r)
	return err
}

func (u *User) Remove() error {
	sb := goqu.From(userTable).
		Delete().
		Where(goqu.Ex{
			userCId: u.Id,
		})
	_, err := ExecQ(u, sb)
	return err
}

type Users []User

func (u Users) GetDb() *sqlx.DB {
	return nil
}

func (u *Users) All() error {
	sb := goqu.From(userTable).
		Select(userCAll)
	return SelectQ(u, sb)
}
func (u *Users) AllT(tx *sqlx.Tx) error {
	sb := goqu.From(userTable).
		Select(userCAll)
	return SelectQ(u, sb, tx)
}
