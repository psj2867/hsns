package models

import (
	"time"

	"github.com/guregu/null/v5"
)

// const (
// 	contentTable     = "content"
// 	contentCAll      = "*"
// 	contentCId       = "id"
// 	contentCName     = "name"
// 	contentCFullName = "fullname"
// )

type Content struct {
	Id       int64
	Content  null.String
	Uuid     string
	CreateAt time.Time
	Uploaded bool
	Images   []Image
}
type Image struct {
	Id        int64
	ContentId int64
	Uuid      string
}

func NewContent() Content {
	return Content{
		CreateAt: time.Now(),
		Uploaded: false,
	}
}
