package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	println("Starting app...")

	// creating a chi router instance
	router := chi.NewRouter()
	// applying the chi middleware
	router.Use(middleware.Logger)

	// defining a chi route
	router.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	// creating a net/http instance
	// using reference of the http.Server even though go implicitly dereferences the instance on each method call
	// why you ask? ask chatgpt
	server := &http.Server{
		Addr: ":3000",
		Handler: router,
	}

	println("Starting server...")
	err := server.ListenAndServe()

	// explicitly panicking because well... it's go
	if err != nil {
		panic(err)
	}
}
