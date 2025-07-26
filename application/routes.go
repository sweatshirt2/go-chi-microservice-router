package application

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sweatshirt2/go-analytics/handler"
	repository "github.com/sweatshirt2/go-analytics/repositories"
)

func (app *App) loadRoutes() {
	// creating a chi router instance
	router := chi.NewRouter()
	// applying the chi middleware
	router.Use(middleware.Logger)

	router.Route("/orders", app.loadOrderRoutes)

	app.router = router
}

func (a *App) loadOrderRoutes(router chi.Router) {
	orderController := handler.OrderController{
		Repo: &repository.OrderRepo{
			Client: a.rds,
		},
	}

	router.Post("/", orderController.Create)
	router.Get("/", orderController.GetAll)
	router.Get("/{id}", orderController.GetById)
	router.Put("/{id}", orderController.Update)
	router.Delete("/{id}", orderController.Delete)
}
