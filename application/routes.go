package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func loadRoutes() *chi.Mux {
	// creating a chi router instance
	router := chi.NewRouter()
	// applying the chi middleware
	router.Use(middleware.Logger)

	// defining a chi route
	router.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	return router
}
