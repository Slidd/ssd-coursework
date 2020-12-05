package main

import (
	"log"
	"ssd-coursework/app"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = app.Init()
	if err != nil {
		log.Fatal(err)
	}

	StartServer()
}

