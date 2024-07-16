package models

import (
	"github.com/jmoiron/sqlx"
)

const (
	cContentTable    = "content"
	cContentUploaded = "uploaded"
)

type Content struct {
	ContentRequest
}

func FromRequestToContent(request ContentRequest) Content {
	return Content{request}
}

func (u *Content) AddT(tx *sqlx.Tx) error {
	sb := from(cContentTable).
		Insert().
		Rows(
			u.ToRecord(),
		)
	return AddQ(sb, &u.Id, tx)
}

func (u *Content) AddImagesT(images Images, tx *sqlx.Tx) error {
	for _, v := range images {
		v.ContentId = u.Id
	}
	return images.AddT(tx)
}
