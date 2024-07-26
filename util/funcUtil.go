package util

import (
	"encoding/json"
	"strings"

	"github.com/sa-/slicefunk"
)

func FlatMap[T any, U any](s []T, f func(T) []U) []U {
	l := slicefunk.Map(s, f)
	return slicefunk.Flatten(l)
}

func GenerateN[T any](f func() T, n int) []T {
	res := make([]T, n)
	for i := range n {
		res[i] = f()
	}
	return res
}

func StringList(s []any) []string {
	return slicefunk.Map(s, func(a any) string {
		return a.(string)
	})
}

func JsonToMap(s []byte) (map[string]any, error) {
	var result map[string]any
	if err := json.Unmarshal(s, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func MustNErrorS(err error, text ...string) {
	if err != nil {
		panic(strings.Join(text, " : ") + " : " + err.Error())
	}
}
