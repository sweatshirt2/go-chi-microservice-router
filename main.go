package main

import (
	"context"

	"github.com/sweatshirt2/go-analytics/application"
)

func main() {
	println("Starting App...")

	app := application.NewApp()
	err := app.Start(context.TODO())

	// explicitly panicking because well... it's go
	if err != nil {
		panic(err)
	}
}
