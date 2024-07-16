package models

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"github.com/thoas/go-funk"
)

type Image struct {
	Id        int64
	ContentId int64 `db:"content_id"`
	Uuid      string
}

const (
	cImageTable     = "image"
	cImageContentId = "content_id"
	cImageUuid      = "uuid"
)

func (u *Image) toRecord() goqu.Record {
	return goqu.Record{
		cImageContentId:      u.ContentId,
		cContentRequestCUuid: u.Uuid,
	}
}

type Images []Image

func (t Images) AddT(tx *sqlx.Tx) error {
	sb := from(cImageTable).
		Insert().
		Rows(
			funk.Map(t, func(a Image) goqu.Record {
				return a.toRecord()
			}),
		)
	_, err := ExecQ(sb, tx)
	return err
}
