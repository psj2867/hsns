package main

import "github.com/psj2867/hsns/server"

func main() {
	s := server.InitServer()
	s.Run(":8080")
}
