package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v5"
)

// const (
// 	contentRequestTable     = "content_request"
// 	contentRequestCAll      = "*"
// 	contentRequestCId       = "id"
// 	contentRequestCName     = "name"
// 	contentRequestCFullName = "fullname"
// )

type ContentRequest struct {
	Id       int64
	UserId   int64
	Content  null.String
	Uuid     string
	CreateAt time.Time
}

func NewContentRequest() ContentRequest {
	return ContentRequest{
		CreateAt: time.Now(),
		Uuid:     uuid.New().String(),
	}
}
