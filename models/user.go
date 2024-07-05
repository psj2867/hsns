package models

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/guregu/null/v5"
)

const (
	table     = "user"
	cAll      = "*"
	cId       = "id"
	cName     = "name"
	cFullName = "fullname"
)

type User struct {
	Id       int64
	Name     string
	Fullname null.String
}

func (u *User) Get(id int) error {

	sb := goqu.From(table).
		Select(cAll).
		Where(goqu.Ex{
			cId: id,
		})
	return GetQ(u, sb)
}

func (u *User) Add() error {
	sb := goqu.From(table).
		Insert().
		Cols(cName, cFullName).
		Vals(goqu.Vals{u.Name, u.Fullname})
	r, err := ExecQ(u, sb)
	if err != nil {
		return err
	}
	err = Refresh(&u.Id, r)
	return err
}

func (u *User) Remove() error {
	sb := goqu.From(table).
		Delete().
		Where(goqu.Ex{
			cId: u.Id,
		})
	_, err := ExecQ(u, sb)
	return err
}

type Users []User

func (u *Users) All() error {
	sb := goqu.From(table).
		Select(cAll)
	return SelectQ(u, sb)
}
