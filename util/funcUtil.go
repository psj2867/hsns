package util

import "github.com/sa-/slicefunk"

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
