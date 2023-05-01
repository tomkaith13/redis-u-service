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

func CfReserve(w http.ResponseWriter, r *http.Request) {
	var client = redisbloom.NewClient(os.Getenv("REDIS_DB_URL"), "nohelp", nil)
	var cfRequest ReserveRequest

	err := json.NewDecoder(r.Body).Decode(&cfRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("req: %+v", cfRequest)
	err = client.Reserve(cfRequest.Name, cfRequest.ErrorRate, cfRequest.Capacity)
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
