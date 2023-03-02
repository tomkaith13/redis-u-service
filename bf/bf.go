package bf

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	redisbloom "github.com/RedisBloom/redisbloom-go"
	redis "github.com/go-redis/redis/v8"
)

type ReserveRequest struct {
	Name      string  `json:"name"`
	ErrorRate float64 `json:"errorRate"`
	Capacity  uint64  `json:"capacity"`
	TtlInSecs uint64  `json:"ttl_in_secs"`
}

type AddItemRequest struct {
	KeyName string `json:"keyName"`
	Item    string `json:"item"`
}

type ExistsRequest struct {
	KeyName string `json:"keyName"`
	Item    string `json:"item"`
}

type DeleteKeyRequest struct {
	KeyName string `json:"keyName"`
}

func BfAddTestFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Println("testing bfadd-test route!!")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("it works!!"))
}

func BfTestSetup(w http.ResponseWriter, r *http.Request) {
	var client = redisbloom.NewClient(os.Getenv("REDIS_DB_URL"), "nohelp", nil)
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

}

func BfReserve(w http.ResponseWriter, r *http.Request) {
	var client = redisbloom.NewClient(os.Getenv("REDIS_DB_URL"), "nohelp", nil)
	var bfRequest ReserveRequest

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

	if bfRequest.TtlInSecs != 0 {
		RedisClient := redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_DB_URL"),
			Password: os.Getenv("REDIS_DB_PASSWORD"),
			DB:       0,
		})
		ctx := context.Background()

		val, err := RedisClient.Expire(ctx, bfRequest.Name, time.Duration(bfRequest.TtlInSecs*uint64(time.Second))).Result()
		if err != nil || !val {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}
	w.WriteHeader(http.StatusCreated)

	b, err := json.Marshal(bfRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)

}

func BfInsert(w http.ResponseWriter, r *http.Request) {
	var client = redisbloom.NewClient(os.Getenv("REDIS_DB_URL"), "nohelp", nil)
	var bfRequest AddItemRequest

	err := json.NewDecoder(r.Body).Decode(&bfRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("req: %+v\n", bfRequest)

	res, err := client.BfInsert(bfRequest.KeyName, 0, 0, 0, true, false, []string{bfRequest.Item})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			w.WriteHeader(http.StatusNotFound)
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

}

func BfExists(w http.ResponseWriter, r *http.Request) {
	var client = redisbloom.NewClient(os.Getenv("REDIS_DB_URL"), "nohelp", nil)
	var bfRequest ExistsRequest

	err := json.NewDecoder(r.Body).Decode(&bfRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("req: %+v\n", bfRequest)
	res, err := client.Exists(bfRequest.KeyName, bfRequest.Item)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if res {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("maybe exists"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("definitely does not exist"))
}

func BfDelete(w http.ResponseWriter, r *http.Request) {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_DB_URL"),
		Password: os.Getenv("REDIS_DB_PASSWORD"),
		DB:       0,
	})
	var ctx = context.Background()
	var bfRequest DeleteKeyRequest

	err := json.NewDecoder(r.Body).Decode(&bfRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	val, err := client.Del(ctx, bfRequest.KeyName).Result()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if val == 0 {
		// this means we dont have anything to delete, hence 404
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("bloomfilter not found!"))
		return
	}

}
