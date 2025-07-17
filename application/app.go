package application

import (
	"context"
	"fmt"
	"net/http"
)

type App struct {
	router http.Handler
}

func NewApp() *App {
	router := loadRoutes()

	app := &App{router: router}

	return app
}

func (app *App) Start(ctx context.Context) error {
	// creating a net/http instance
	// using reference of the http.Server even though go implicitly dereferences the instance on each method call
	server := &http.Server{
		Addr: ":3000",
		Handler: app.router,
	}

	err := server.ListenAndServe()

	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
