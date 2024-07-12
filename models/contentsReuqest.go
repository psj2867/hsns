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
	return WrapGet(u, id,
		func(wgo *wrapSqlOptions) { wgo.tableName = cContentRequestTable },
	)
}
func (u *ContentRequest) GetT(id int64, tx *sqlx.Tx) error {
	return WrapGet(u, id,
		func(wgo *wrapSqlOptions) { wgo.tableName = cContentRequestTable },
		func(wgo *wrapSqlOptions) { wgo.tx = tx },
	)
}

func (u *ContentRequest) Add() error {
	sb := from(cContentRequestTable).
		Insert().
		Rows(
			goqu.Record{
				cContentRequestCUserId:   u.UserId,
				cContentRequestCContent:  u.Content,
				cContentRequestCUuid:     u.Uuid,
				cContentRequestCCreateAt: u.CreateAt,
			},
		)
	return AddQ(u, sb, &u.Id)
}
func (u *ContentRequest) RemoveT(tx *sqlx.Tx) error {
	return WrapRemove(u, u.Id,
		func(wgo *wrapSqlOptions) { wgo.tableName = cContentRequestTable },
		func(wgo *wrapSqlOptions) { wgo.tx = tx },
	)
}
