package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tomkaith13/redis-u-service/bf"
)

func main() {
	fmt.Println("testing!!")
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/bfadd-test", bf.BfAddTestFunc)

	r.Post("/bf-test-setup", bf.BfTestSetup)

	r.Post("/bf-reserve", bf.BfReserve)

	r.Post("/bf-insert", bf.BfInsert)

	r.Delete("/bf", bf.BfDelete)

	r.Get("/bf-exists", bf.BfExists)

	log.Fatal(http.ListenAndServe(":8080", r))

}
