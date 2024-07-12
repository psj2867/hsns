package models

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
)

const (
	cContentTable    = "content"
	cContentUploaded = "uploaded"
)

type Content struct {
	ContentRequest
	Uploaded bool
}

func FromRequestToContent(request ContentRequest) Content {
	return Content{
		ContentRequest: request,
		Uploaded:       true,
	}
}
func (u *Content) toRecord() goqu.Record {
	return goqu.Record{
		cContentRequestCUserId:   u.UserId,
		cContentRequestCContent:  u.Content,
		cContentRequestCUuid:     u.Uuid,
		cContentRequestCCreateAt: u.CreateAt,
		cContentUploaded:         u.Uploaded,
	}
}

func (u *Content) AddT(tx *sqlx.Tx) error {
	sb := from(cContentTable).
		Insert().
		Rows(
			u.toRecord(),
		)
	return AddQ(u, sb, &u.Id, tx)
}

func (u *Content) AddImagesT(images Images, tx *sqlx.Tx) error {
	return images.AddT(tx)
}
