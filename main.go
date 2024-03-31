package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func main() {

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/", handleGet)

	fmt.Println("Server running on port 3000")

	http.ListenAndServe(":3000", router)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, api world!"))
}
