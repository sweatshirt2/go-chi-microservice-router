package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type App struct {
	router http.Handler
	rds *redis.Client
}

func NewApp() *App {
	rds := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB: 9,
	})

	app := &App{
		rds: rds,
	}
	app.loadRoutes()

	return app
}

func (app *App) Start(ctx context.Context) error {
	// creating a net/http instance
	// using reference of the http.Server even though go implicitly dereferences the instance on each method call
	server := &http.Server{
		Addr: ":3000",
		Handler: app.router,
	}

	err := app.rds.Ping(ctx).Err()

	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	defer func ()  {
		if close_err := app.rds.Close(); close_err != nil {
			fmt.Println("failed to close redis", err)
		}
	}()

	println("Connected to redis successfully!\nStarting server...")

	ch := make(chan error, 1)

	go func ()  {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}

		close(ch)	
	}()

	// ch_err, open := <-ch
	ch_err := <-ch

	select {
	case ch_err = <-ch:
		return ch_err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		return server.Shutdown(timeout)
	}
}
