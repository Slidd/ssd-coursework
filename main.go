package main

import (
	"log"
	"ssd-coursework/app"
)

func main() {
	err := app.Init()
	if err != nil {
		log.Fatal(err)
	}
	StartServer()
}
