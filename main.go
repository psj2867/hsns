package main

import "github.com/psj2867/hsns/router"

func main() {
	s := router.InitServer()
	s.Run(":8080")
}
