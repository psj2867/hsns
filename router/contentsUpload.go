package router

import (
	"context"
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/psj2867/hsns/config"
	"github.com/psj2867/hsns/models"
	"github.com/sa-/slicefunk"
)

type uploadRequest struct {
	Content string `form:"content" binding:"required"`
	Images  int    `form:"images" binding:"-"`
}

func (t *uploadRequest) toContentRequest(userId int64) *models.ContentRequest {
	return &models.ContentRequest{
		UserId:   userId,
		Content:  t.Content,
		CreateAt: time.Now(),
		Uuid:     uuid.New().String(),
	}
}

func (t *contents) upload(c *gin.Context) {
	var ur uploadRequest
	if err := c.Bind(&ur); err != nil {
		c.JSON(400, err.Error())
		return
	}

	user := models.GetUserInfoInContext(c)
	contentsReuqest := ur.toContentRequest(user.Id)
	if err := contentsReuqest.Add(); err != nil {
		c.JSON(403, err.Error())
		return
	}

	uploadToken := config.CreateUploadToken(&config.UploadTokenInfo{
		Uuid:     contentsReuqest.Uuid,
		Images:   ur.Images,
		CreateAt: contentsReuqest.CreateAt,
	})
	c.Status(200)
	c.JSON(200, map[string]any{
		"uuids": uploadToken.Images,
		"token": uploadToken.String(),
	})
}

type uploadTokenRequest struct {
	Token string `form:"token"`
}

func (t *contents) uploadSuccess(c *gin.Context) {

	var request uploadTokenRequest
	if err := c.Bind(&request); err != nil {
		return
	}
	token := request.Token
	returnToken := config.ParseReturnToken(token)
	// start tx
	tx, _ := config.GetDb().BeginTxx(context.Background(), &sql.TxOptions{})
	defer tx.Commit()

	var contentRequest models.ContentRequest
	contentRequest.GetByUuidT(returnToken.Uuid, tx)
	contentRequest.RemoveT(tx)

	content := models.FromRequestToContent(contentRequest)
	content.AddT(tx)
	content.AddImagesT(toImages(returnToken.GetImages()), tx)
}
func toImages(imageUuids []string) models.Images {
	return slicefunk.Map(imageUuids, func(uuid string) models.Image {
		return models.Image{
			Uuid: uuid,
		}
	})
}

func (t *contents) uploadFail(c *gin.Context) {
	var request uploadTokenRequest
	if err := c.Bind(&request); err != nil {
		return
	}
	token := request.Token
	parsedToken := config.UploadToken{}
	parsedToken.Parse(token)
	contentRequest := models.ContentRequest{Uuid: parsedToken.GetUuid()}
	contentRequest.RemoveByUuid()
}
func IterateItems(yield func(int) bool) {
	items := []int{1, 2, 3}
	for _, v := range items {
		if !yield(v) {
			return
		}
	}
}
