package service

import (
	"context"
	"github.com/inasknh/simple-poke-app/internal/api"
	"github.com/inasknh/simple-poke-app/internal/model"
	"github.com/inasknh/simple-poke-app/internal/repository"
)

type service struct {
	dbRepository    repository.Repository
	redisRepository repository.RedisRepository
	client          api.Client
}

type Service interface {
	SyncData(ctx context.Context) error
	GetItems(ctx context.Context) (*model.BerriesResponse, error)
}

func NewService(repository repository.Repository,
	redisRepository repository.RedisRepository,
	client api.Client) Service {
	return &service{
		dbRepository:    repository,
		client:          client,
		redisRepository: redisRepository,
	}
}

func (s *service) SyncData(ctx context.Context) error {
	// get data from client
	res, err := s.client.GetBerries(ctx, api.BerriesRequest{})
	if err != nil {
		return err
	}

	// insert to db
	berries := constructBerries(res)
	err = s.dbRepository.CreateBerry(ctx, berries)
	if err != nil {
		return err
	}

	return nil
}

func constructBerries(res *api.BerriesResponse) []model.Berry {
	berries := make([]model.Berry, 0, len(res.Results))
	for _, result := range res.Results {
		berries = append(berries, model.Berry{
			Name: result.Name,
			URL:  result.Url,
		})
	}

	return berries
}

func (s *service) GetItems(ctx context.Context) (*model.BerriesResponse, error) {

	cacheRes, err := s.redisRepository.GetData(ctx)
	if cacheRes != nil {
		return cacheRes, nil
	}

	data, err := s.dbRepository.FetchBerries(ctx)
	if err != nil {
		return nil, err
	}

	berries := make([]model.Berry, 0, len(data.Berries))
	for _, berry := range data.Berries {
		berries = append(berries, model.Berry{
			Name: berry.Name,
			URL:  berry.URL,
		})
	}

	response := &model.BerriesResponse{Berries: berries}

	err = s.redisRepository.SetData(ctx, response)

	// regardless the return from SetData, it should be return response
	return response, nil
}
