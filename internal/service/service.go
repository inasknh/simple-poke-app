package service

import (
	"context"
	"github.com/inasknh/simple-poke-app/internal/api"
	"github.com/inasknh/simple-poke-app/internal/model"
	"github.com/inasknh/simple-poke-app/internal/repository"
)

type service struct {
	repository repository.Repository
	client     api.Client
}

type Service interface {
	SyncData(ctx context.Context) error
	GetItems(ctx context.Context) (*model.BerriesResponse, error)
}

func NewService(repository repository.Repository, client api.Client) Service {
	return &service{
		repository: repository,
		client:     client,
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
	err = s.repository.CreateBerry(ctx, berries)
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

}
