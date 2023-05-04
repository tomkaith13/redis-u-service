package cf

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	redisbloom "github.com/RedisBloom/redisbloom-go"
	"github.com/go-redis/redis/v8"
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

type DelItemRequest struct {
	KeyName string `json:"keyName"`
	Item    string `json:"item"`
}

func CfReserve(w http.ResponseWriter, r *http.Request) {
	var client = redisbloom.NewClient(os.Getenv("REDIS_DB_URL"), "nohelp", nil)
	var cfRequest ReserveRequest

	err := json.NewDecoder(r.Body).Decode(&cfRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("req: %+v", cfRequest)
	_, err = client.CfReserve(cfRequest.Name, int64(cfRequest.Capacity), 0, 0, 0)
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

	if cfRequest.TtlInSecs != 0 {
		RedisClient := redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_DB_URL"),
			Password: os.Getenv("REDIS_DB_PASSWORD"),
			DB:       0,
		})
		ctx := context.Background()

		val, err := RedisClient.Expire(ctx, cfRequest.Name, time.Duration(cfRequest.TtlInSecs*uint64(time.Second))).Result()
		if err != nil || !val {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}
	w.WriteHeader(http.StatusCreated)

	b, err := json.Marshal(cfRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func CfInsert(w http.ResponseWriter, r *http.Request) {
	var client = redisbloom.NewClient(os.Getenv("REDIS_DB_URL"), "nohelp", nil)
	var cfRequest AddItemRequest

	err := json.NewDecoder(r.Body).Decode(&cfRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("req: %+v\n", cfRequest)

	res, err := client.CfAdd(cfRequest.KeyName, cfRequest.Item)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("err: No CF with keyName exists. Use POST /cf-reserve to create a new one"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("err: " + err.Error()))
		return
	}
	if !res {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("item did not get added. It could already exist. Check /cf-exists"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("added"))

}

func CfExists(w http.ResponseWriter, r *http.Request) {
	var client = redisbloom.NewClient(os.Getenv("REDIS_DB_URL"), "nohelp", nil)
	var cfRequest ExistsRequest

	err := json.NewDecoder(r.Body).Decode(&cfRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("req: %+v\n", cfRequest)
	res, err := client.CfExists(cfRequest.KeyName, cfRequest.Item)

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

func CfDeleteItem(w http.ResponseWriter, r *http.Request) {
	var client = redisbloom.NewClient(os.Getenv("REDIS_DB_URL"), "nohelp", nil)
	var cfRequest DelItemRequest

	err := json.NewDecoder(r.Body).Decode(&cfRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("req: %+v\n", cfRequest)
	res, err := client.CfDel(cfRequest.KeyName, cfRequest.Item)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if !res {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("item not found"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Deleted item"))
}
