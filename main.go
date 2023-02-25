package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	fmt.Println("testing!!")
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/bfadd-test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("testing bfadd-test route!!")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("it works!!"))
	})

	log.Fatal(http.ListenAndServe(":8080", r))

}
