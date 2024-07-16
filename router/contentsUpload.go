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
	"github.com/psj2867/hsns/util"
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

	uploadToken := createUploadToken(contentsReuqest, &ur).String()
	c.Status(200)
	c.Writer.WriteString(uploadToken)
}
func gerateImageUuid() string {
	return uuid.New().String()
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
	returnToken := parseReturnToken(token)
	// start tx
	tx, _ := config.GetDb().BeginTxx(context.Background(), &sql.TxOptions{})
	defer tx.Commit()

	var contentRequest models.ContentRequest
	contentRequest.GetByUuidT(returnToken.Uuid, tx)
	contentRequest.RemoveT(tx)

	content := models.FromRequestToContent(contentRequest)
	content.AddT(tx)
	content.AddImagesT(returnToken.getImages(), tx)
}

func (t *contents) uploadFail(c *gin.Context) {
	var request uploadTokenRequest
	if err := c.Bind(&request); err != nil {
		return
	}
	token := request.Token
	parsedToken := uploadToken{}
	parsedToken.parse(token)
	contentRequest := models.ContentRequest{Uuid: parsedToken.getUuid()}
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

type uploadToken struct {
	data   map[string]interface{}
	images []string
}

func createUploadToken(cr *models.ContentRequest, req *uploadRequest) *uploadToken {
	m := uploadToken{data: map[string]interface{}{}}
	m.images = util.GenerateN(gerateImageUuid, req.Images)
	m.data["createAt"] = cr.CreateAt
	m.data["uuid"] = cr.Uuid
	m.data["images"] = m.images
	return &m
}

func (t *uploadToken) String() string {
	jr, _ := json.Marshal(t.data)
	encoded, _ := config.UploadTokenEnDecoder.Encode(jr)
	return string(encoded)
}
func (t *uploadToken) parse(data string) error {
	decodedData, _ := config.UploadTokenEnDecoder.Decode([]byte(data))
	err := json.Unmarshal(decodedData, &t.data)
	if err != nil {
		return err
	}
	t.images = t.data["images"].([]string)
	_ = t.data["uuid"].(string)
	return nil
}
func (t *uploadToken) getUuid() string {
	uuid := t.data["uuid"]
	if uuid == nil {
		return ""
	}
	return uuid.(string)
}

type returnToken struct {
	Uuid           string   `json:"uuid"`
	RequestImages  []string `json:"requestImages"`
	UploadedImages []string `json:"uploadedImages"`
}

func (t *returnToken) getImages() models.Images {
	return slicefunk.Map(t.UploadedImages, func(uuid string) models.Image {
		return models.Image{
			Uuid: uuid,
		}
	})
}

func parseReturnToken(token string) *returnToken {
	decoded, _ := config.ReturnTokenEnDecoder.Decode([]byte(token))
	var r returnToken
	err := json.Unmarshal(decoded, &r)
	if err != nil {
		return nil
	}
	return &r
}
func MergeTokens(tokens ...*returnToken) *returnToken {
	if len(tokens) == 0 {
		return nil
	}
	base := tokens[0]
	base.UploadedImages = util.FlatMap(tokens, func(t *returnToken) []string {
		return t.UploadedImages
	})
	return base
}
