package router

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/psj2867/hsns/config"
	"github.com/psj2867/hsns/models"
)

type uploadRequest struct {
	Content string `form:"content" binding:"required"`
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
	user := models.GetUserInfoInContext(c)
	var ur uploadRequest
	c.Bind(&ur)
	contentsReuqest := ur.toContentRequest(user.Id)
	if err := contentsReuqest.Add(); err != nil {
		c.JSON(403, err.Error())
		return
	}
	uploadToken := t.createUploadToken(contentsReuqest).String()
	c.JSON(200, uploadToken)
}
func (t *contents) createUploadToken(cr *models.ContentRequest) *uploadToken {
	m := uploadToken{data: map[string]interface{}{}}
	m.data["user"] = cr.UserId
	m.data["CreateAt"] = cr.CreateAt
	return &m
}

type uploadSuccessRequest struct {
	Token string `form:"token"`
}

func (t *contents) uploadSuccess(c *gin.Context) {

	var request uploadSuccessRequest
	if err := c.Bind(&request); err != nil {
		return
	}
	token := request.Token
	returnToken := parseReturnToekn(token)
	returnToken.shouldHave()
	// start tx
	tx, _ := config.GetDb().BeginTxx(context.Background(), &sql.TxOptions{})
	defer tx.Commit()

	var contentRequest models.ContentRequest
	contentRequest.GetT(returnToken.get("").(int64), tx)
	contentRequest.RemoveT(tx)

	content := models.FromRequestToContent(contentRequest)
	content.AddT(tx)
	content.AddImagesT(returnToken.get("images").(models.Images), tx)
}

func (t *contents) uploadFail(c *gin.Context) {
	// 토큰 확인
	// 삭제
}

type uploadToken struct {
	data map[string]interface{}
}

func (t *uploadToken) String() string {
	jr, _ := json.Marshal(t.data)
	encoded := t.encode(jr)
	return encoded
}
func (t *uploadToken) encode(data []byte) string {
	bEnc, _ := config.UploadTokenEnDecoder{}.Encode(data)
	return string(bEnc)
}

type returnToken struct {
	data map[string]interface{}
}

func parseReturnToekn(token string) *returnToken {
	decoded, _ := config.ReturnTokenEnDecoder{}.Decode([]byte(token))
	var r returnToken
	_ = token
	// TODO
	return &r
}
func (t *returnToken) shouldHave(keys ...string) error {
	return nil
}
func (t *returnToken) get(key string) interface{} {
	v, _ := t.data[key]
	return v
}
