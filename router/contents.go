package router

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/guregu/null/v5"
	"github.com/psj2867/hsns/config"
	"github.com/psj2867/hsns/models"
)

type contents struct{}

func (t contents) setRouter(group *gin.RouterGroup) {
	group.GET("/", t.getContent)
	group.POST("/upload", t.upload)
	group.GET("/upload/success", t.uploadSuccess)
	group.GET("/upload/fail", t.uploadFail)

}
func (t *contents) getContent(c *gin.Context) {
}

type uploadRequest struct {
	Content string `form:"content" binding:"required"`
}

func (t *uploadRequest) to(userId int64) *models.ContentRequest {
	return &models.ContentRequest{
		UserId:   userId,
		Content:  null.StringFrom(t.Content),
		CreateAt: time.Now(),
		Uuid:     uuid.New().String(),
	}
}
func (t *contents) upload(c *gin.Context) {
	userId, _ := c.Request.Context().Value("userid").(int64)
	var ur uploadRequest
	c.Bind(&ur)
	contentsReuqest := ur.to(userId)
	uploadToken := t.createUploadToken(contentsReuqest)
	c.JSON(200, uploadToken)
}
func (t *contents) createUploadToken(cr *models.ContentRequest) string {
	m := make(map[string]interface{})
	m["user"] = cr.UserId
	m["CreateAt"] = cr.CreateAt
	jr, _ := json.Marshal(m)
	bEnc := base64.StdEncoding.EncodeToString(jr)
	return bEnc
}

func (t *contents) uploadSuccess(c *gin.Context) {
	// user, ok := c.Request.Context().Value("user").(string)
	// if ok == false {
	// 	return
	// }
	// id, _ := c.Params.Get("id")
	// token, _ := c.Params.Get("token")
	// images := parseToken(token)
	// start tx
	tx, _ := config.GetDb().BeginTxx(context.Background(), &sql.TxOptions{})
	defer tx.Commit()
	_ = tx
	// contentRequest := models.ContentRequest.Get(id, tx)
	// contentRequest.remove(tx)

	// content := models.content{}.from(contentRequest)
	// content.Add(tx)
	// content.AddImages(token, tx)
	// end tx
}
func (t *contents) uploadFail(c *gin.Context) {
	// 토큰 확인
	// 삭제
}

func getToken(c *gin.Context) string {
	return ""
}
func setToken(c *gin.Context) {

}
