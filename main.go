package main

import (
	"fmt"

	"github.com/psj2867/hsns/config"
)

func main() {
	// s := server.InitServer()
	// defer server.DeferServer(s)
	// s.Run(":8080")
	f1()
}
func f1() {
	a, _ := config.UploadTokenEnDecoder{}.Encode([]byte("asdf"))
	b, _ := config.UploadTokenEnDecoder{}.Decode(a)
	fmt.Printf("a: %v\n", string(a))
	fmt.Printf("b: %v\n", string(b))
}
