package repository

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redismock/v7"
	"github.com/inasknh/simple-poke-app/internal/config"
	"github.com/inasknh/simple-poke-app/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_redisRepository_GetData(t *testing.T) {
	rd, mock := redismock.NewClientMock()
	defer rd.Close()

	cfg := config.Configurations{}
	repo := NewRedisRepository(rd, cfg)

	ctx := context.Background()
	expected := &model.BerriesResponse{
		Berries: []model.Berry{
			{
				Name: "1",
				URL:  "1",
			},
		},
	}

	data, _ := json.Marshal(expected)

	mock.ExpectGet("items").SetVal(string(data))

	result, err := repo.GetData(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func Test_redisRepository_GetData_NotFound(t *testing.T) {
	rd, mock := redismock.NewClientMock()
	defer rd.Close()

	cfg := config.Configurations{}
	repo := NewRedisRepository(rd, cfg)

	ctx := context.Background()

	mock.ExpectGet("items").RedisNil()

	result, err := repo.GetData(ctx)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func Test_redisRepository_GetData_UnmarshalError(t *testing.T) {
	rd, mock := redismock.NewClientMock()
	defer rd.Close()

	cfg := config.Configurations{}
	repo := NewRedisRepository(rd, cfg)

	ctx := context.Background()

	mock.ExpectGet("items").SetVal("invalid-json")

	result, err := repo.GetData(ctx)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid character")
}

func Test_redisRepository_SetData(t *testing.T) {
	rd, mock := redismock.NewClientMock()
	defer rd.Close()

	cfg := config.Configurations{
		App: config.AppConfiguration{
			TTL: 5,
		},
	}
	repo := NewRedisRepository(rd, cfg)

	ctx := context.Background()
	expected := &model.BerriesResponse{Berries: []model.Berry{
		{
			Name: "1",
			URL:  "1",
		},
	}}

	data, _ := json.Marshal(expected)

	mock.ExpectSet("items", data, time.Duration(5)*time.Minute).SetVal("OK")
	err := repo.SetData(ctx, expected)
	assert.NoError(t, err)

}

func Test_redisRepository_SetData_Failure(t *testing.T) {
	rd, mock := redismock.NewClientMock()
	defer rd.Close()

	cfg := config.Configurations{
		App: config.AppConfiguration{
			TTL: 5,
		},
	}
	repo := NewRedisRepository(rd, cfg)

	ctx := context.Background()
	expected := &model.BerriesResponse{Berries: []model.Berry{
		{
			Name: "1",
			URL:  "1",
		},
	}}

	data, _ := json.Marshal(expected)

	mock.ExpectSet("items", data, time.Duration(5)*time.Minute).
		SetErr(errors.New("an error"))
	err := repo.SetData(ctx, expected)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "an error")

}
