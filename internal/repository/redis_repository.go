package repository

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v7"
	"github.com/inasknh/simple-poke-app/internal/config"
	"github.com/inasknh/simple-poke-app/internal/model"
	"time"
)

type redisRepository struct {
	cache  *redis.Client
	config config.Configurations
}

func NewRedisRepository(cache *redis.Client, config config.Configurations) RedisRepository {
	return &redisRepository{cache: cache, config: config}
}

type RedisRepository interface {
	GetData(ctx context.Context) (*model.BerriesResponse, error)
	SetData(ctx context.Context, response *model.BerriesResponse) error
}

func (r *redisRepository) GetData(ctx context.Context) (*model.BerriesResponse, error) {
	res, err := r.cache.Get("items").Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	var berries model.BerriesResponse
	err = json.Unmarshal(res, &berries)
	if err != nil {
		return nil, err
	}

	return &berries, nil
}

func (r *redisRepository) SetData(ctx context.Context, response *model.BerriesResponse) error {
	data, err := json.Marshal(response)
	if err != nil {
		return err
	}

	_, err = r.cache.Set("items", data, time.Duration(r.config.App.TTL)*time.Minute).Result()
	if err != nil {
		return err
	}

	return nil
}
