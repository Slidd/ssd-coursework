package main

import (
	"log"
	"ssd-coursework/app"
)

// entrypoint
func main() {
	err := app.Init()
	if err != nil {
		log.Fatal(err)
	}
	StartServer()
}
