package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	println("Starting app...")

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	server := &http.Server{
		Addr: ":3000",
		Handler: router,
	}

	println("Starting server...")
	err := server.ListenAndServe()

	if err != nil {
		panic(err)
	}
}
