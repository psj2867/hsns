package models

import (
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
)

const (
	cContentRequestTable     = "content_request"
	cContentRequestCId       = "id"
	cContentRequestCUserId   = "user_id"
	cContentRequestCContent  = "content"
	cContentRequestCUuid     = "uuid"
	cContentRequestCCreateAt = "create_at"
)

type ContentRequest struct {
	Id       int64
	UserId   int64 `db:"user_id"`
	Content  string
	Uuid     string
	CreateAt time.Time `db:"create_at"`
}

func (u *ContentRequest) Get(id int64) error {
	return wrapGet(u, id, withTable(cContentRequestTable))
}
func (u *ContentRequest) GetByUuid(uuid string) error {
	return wrapGet(u, uuid, withTable(cContentRequestTable),
		func(wgo *wrapSqlOptions) { wgo.idName = cContentRequestCUuid },
	)
}
func (u *ContentRequest) GetByUuidT(uuid string, tx *sqlx.Tx) error {
	return wrapGet(u, uuid, withTable(cContentRequestTable), withTx(tx),
		func(wgo *wrapSqlOptions) { wgo.idName = cContentRequestCUuid },
	)
}

func (u *ContentRequest) Add() error {
	sb := from(cContentRequestTable).
		Insert().
		Rows(u.ToRecord())
	return AddQ(sb, &u.Id)
}

func (u *ContentRequest) ToRecord() *goqu.Record {
	return &goqu.Record{
		cContentRequestCUserId:   u.UserId,
		cContentRequestCContent:  u.Content,
		cContentRequestCUuid:     u.Uuid,
		cContentRequestCCreateAt: u.CreateAt,
	}
}
func (u *ContentRequest) RemoveByUuid() error {
	return wrapRemove(u.Uuid, withTable(cContentRequestTable),
		func(wgo *wrapSqlOptions) { wgo.idName = cContentRequestCUuid },
	)
}

func (u *ContentRequest) RemoveT(tx *sqlx.Tx) error {
	return wrapRemove(u.Id, withTable(cContentRequestTable), withTx(tx))
}
