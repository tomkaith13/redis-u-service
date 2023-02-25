package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	redisbloom "github.com/RedisBloom/redisbloom-go"
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

	r.Post("/bfadd-test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("testing bfadd-test route!!")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("it works!!"))
	})

	r.Post("/bf-test-setup", func(w http.ResponseWriter, r *http.Request) {
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

	r.Post("/bf-reserve", func(w http.ResponseWriter, r *http.Request) {
		var client = redisbloom.NewClient("redis-server:6379", "nohelp", nil)
		var bfRequest bf.ReserveRequest

		err := json.NewDecoder(r.Body).Decode(&bfRequest)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Printf("req: %+v", bfRequest)

		err = client.Reserve(bfRequest.Name, bfRequest.ErrorRate, bfRequest.Capacity)
		if err != nil {
			if strings.Contains(err.Error(), "item exists") {
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte("err: " + err.Error()))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("err: " + err.Error()))
			return
		}
		w.WriteHeader(http.StatusCreated)

		b, err := json.Marshal(bfRequest)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(b)

	})

	r.Post("/bf-insert", func(w http.ResponseWriter, r *http.Request) {
		var client = redisbloom.NewClient("redis-server:6379", "nohelp", nil)
		var bfRequest bf.AddItemRequest

		err := json.NewDecoder(r.Body).Decode(&bfRequest)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Printf("req: %+v\n", bfRequest)

		res, err := client.BfInsert(bfRequest.KeyName, 0, 0, 0, true, false, []string{bfRequest.Item})
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("err: No BF with keyName exists. Use POST /bf-reserve to create a new one"))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("err: " + err.Error()))
			return
		}

		for _, r := range res {
			fmt.Println("res: ", r)
			if r == 0 {
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte("item may already exist"))
				return
			}
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("added"))

	})

	log.Fatal(http.ListenAndServe(":8080", r))

}
