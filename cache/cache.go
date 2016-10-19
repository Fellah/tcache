package cache

import (
	"gopkg.in/redis.v4"
	"os"
	"strconv"
	"github.com/fellah/tcache/log"
)

var redis_client *redis.Client

func Init() {
	redis_addr := os.Getenv("REDIS_ADDR")
	redis_pass := os.Getenv("REDIS_PASS")
	redis_db   := os.Getenv("REDIS_DB")

	if redis_addr == "" {
		redis_addr = "localhost:6379"
	}

	if redis_db == "" {
		redis_db = "0"
	}

	redis_db_num, err := strconv.Atoi(redis_db)

	if err != nil {
		log.Error.Fatalln(err)
	}

	redis_client = redis.NewClient(&redis.Options{
		Addr:     redis_addr,
		Password: redis_pass,
		DB:       redis_db_num,
	})
}
