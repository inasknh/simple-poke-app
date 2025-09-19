package cache

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/inasknh/simple-poke-app/internal/config"
	"log"
)

func NewRedis(config config.Cache) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password: config.Password,
	})

	if err := rdb.Ping().Err(); err != nil {
		log.Fatalf("Redis cannot be pinged %s", err)
	}

	return rdb
}
