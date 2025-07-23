package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/sweatshirt2/go-analytics/application"
)

func main() {
	app := application.NewApp()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// err := app.Start(context.TODO())
	err := app.Start(ctx)

	// explicitly panicking because well... it's go
	if err != nil {
		panic(err)
	}
}
