package config

import (
	"encoding/base64"
)

type EnDeocoder interface {
	Encode([]byte) ([]byte, error)
	Decode([]byte) ([]byte, error)
}

var UploadTokenEnDecoder uploadTokenEnDecoder

type uploadTokenEnDecoder struct{}

func (t uploadTokenEnDecoder) Encode(src []byte) ([]byte, error) {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	return dst, nil
}
func (t uploadTokenEnDecoder) Decode(src []byte) ([]byte, error) {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	n, err := base64.StdEncoding.Decode(dst, []byte(src))
	return dst[:n], err
}

var ReturnTokenEnDecoder returnTokenEnDecoder

type returnTokenEnDecoder struct{}

func (t returnTokenEnDecoder) Encode(src []byte) ([]byte, error) {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	return dst, nil
}
func (t returnTokenEnDecoder) Decode(src []byte) ([]byte, error) {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	n, err := base64.StdEncoding.Decode(dst, []byte(src))
	return dst[:n], err
}

func GetJwtSecretKey() []byte {
	return []byte("secretkey")
}
