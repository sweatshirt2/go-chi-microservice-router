package application

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sweatshirt2/go-analytics/handler"
)

func loadRoutes() *chi.Mux {
	// creating a chi router instance
	router := chi.NewRouter()
	// applying the chi middleware
	router.Use(middleware.Logger)

	router.Route("/orders", loadOrderRoutes)

	return router
}

func loadOrderRoutes(router chi.Router) {
	orderController := handler.OrderController{}

	router.Post("/", orderController.Create)
	router.Get("/", orderController.GetAll)
	router.Get("/{id}", orderController.GetById)
	router.Put("/{id}", orderController.Update)
	router.Delete("/{id}", orderController.Delete)
}
