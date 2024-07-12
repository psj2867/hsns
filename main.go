package main

import "github.com/psj2867/hsns/server"

func main() {
	s := server.InitServer()
	defer server.DeferServer(s)
	s.Run(":8080")
	// f1()
}
func f1() {

}
