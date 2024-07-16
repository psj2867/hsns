package config

import "encoding/base64"

type EnDeocoder interface {
	Encode([]byte) ([]byte, error)
	Decode([]byte) ([]byte, error)
}

type UploadTokenEnDecoder struct{}

func (t UploadTokenEnDecoder) Encode(src []byte) ([]byte, error) {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	return dst, nil
}
func (t UploadTokenEnDecoder) Decode(src []byte) ([]byte, error) {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	n, err := base64.StdEncoding.Decode(dst, []byte(src))
	return dst[:n], err
}

type ReturnTokenEnDecoder struct{}

func (t ReturnTokenEnDecoder) Encode(src []byte) ([]byte, error) {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	return dst, nil
}
func (t ReturnTokenEnDecoder) Decode(src []byte) ([]byte, error) {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	n, err := base64.StdEncoding.Decode(dst, []byte(src))
	return dst[:n], err
}
