package config

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/psj2867/hsns/util"
)

type UploadToken struct {
	data   map[string]interface{}
	Images []string
}
type UploadTokenInfo struct {
	Images   int
	Uuid     string
	CreateAt time.Time
}

func CreateUploadToken(info *UploadTokenInfo) *UploadToken {
	m := UploadToken{data: map[string]interface{}{}}
	m.Images = util.GenerateN(gerateImageUuid, info.Images)
	m.data["createAt"] = info.CreateAt
	m.data["uuid"] = info.Uuid
	m.data["imageUuids"] = m.Images
	return &m
}

func gerateImageUuid() string {
	return uuid.New().String()
}

func (t *UploadToken) String() string {
	jr, _ := json.Marshal(t.data)
	encoded, _ := UploadTokenEnDecoder.Encode(jr)
	return string(encoded)
}
func (t *UploadToken) Parse(data string) error {
	decodedData, _ := UploadTokenEnDecoder.Decode([]byte(data))
	err := json.Unmarshal(decodedData, &t.data)
	if err != nil {
		return err
	}
	t.Images = t.data["images"].([]string)
	_ = t.data["uuid"].(string)
	return nil
}
func (t *UploadToken) GetUuid() string {
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
type ReturnTokenInfo struct {
	returnToken
}

func CreateReturnToken(info ReturnTokenInfo) *returnToken {
	return &info.returnToken
}

func (t *returnToken) GetImages() []string {
	return t.UploadedImages
}

func ParseReturnToken(token string) *returnToken {
	decoded, _ := ReturnTokenEnDecoder.Decode([]byte(token))
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
