package main

import (
	"fmt"
	"log"
	"net/http"

	redisbloom "github.com/RedisBloom/redisbloom-go"
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

	r.Post("/bfadd-setup", func(w http.ResponseWriter, r *http.Request) {
		var client = redisbloom.NewClient("redis-server:6379", "nohelp", nil)
		res, err := client.Add("testBF", "works")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error returned by redisbloom pkg.." + err.Error()))
			return
		}

		if res {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("item doesnt exist .... new item added!!"))
			return
		}

		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("item maybe exist!!"))

	})

	log.Fatal(http.ListenAndServe(":8080", r))

}
