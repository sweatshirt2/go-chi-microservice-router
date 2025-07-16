package main

import "net/http"

func main() {
	println("Starting app...")

	server := &http.Server{
		Addr: ":3000",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello World"))
		}),
	}

	println("Starting server...")
	err := server.ListenAndServe()

	if err != nil {
		panic(err)
	}
}
