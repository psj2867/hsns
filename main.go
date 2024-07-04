package main

import "github.com/psj2867/hsns/config"

func main() {
	r := config.InitServer()
	if err := r.Run(":8080"); err != nil {
		return
	}
}
