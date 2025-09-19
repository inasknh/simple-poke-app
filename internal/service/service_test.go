package service

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v7"
	"github.com/inasknh/simple-poke-app/internal/api"
	mocks2 "github.com/inasknh/simple-poke-app/internal/mocks/api"
	mocks "github.com/inasknh/simple-poke-app/internal/mocks/repository"
	"github.com/inasknh/simple-poke-app/internal/model"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func Test_service_SyncData(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		mockFunc func() *service
	}{
		{
			name: "given an error when getBerries from client should return an error",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
			mockFunc: func() *service {
				mockDB := &mocks.Repository{}
				mockRedis := &mocks.RedisRepository{}
				mockClient := &mocks2.Client{}

				mockClient.
					On("GetBerries", mock.Anything, api.BerriesRequest{}).
					Return(nil, errors.New("an error"))
				return &service{
					dbRepository:    mockDB,
					redisRepository: mockRedis,
					client:          mockClient,
				}
			},
		},
		{
			name: "given an error when CreateBerry to database should return an error",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
			mockFunc: func() *service {
				mockDB := &mocks.Repository{}
				mockRedis := &mocks.RedisRepository{}
				mockClient := &mocks2.Client{}

				mockClient.
					On("GetBerries", mock.Anything, api.BerriesRequest{}).
					Return(&api.BerriesResponse{
						Count:    0,
						Next:     "",
						Previous: "",
						Results: []api.Berry{
							{
								Name: "1",
								Url:  "1",
							},
						},
					}, nil)

				mockDB.On("CreateBerry", mock.Anything, []model.Berry{
					{
						Name: "1",
						URL:  "1",
					},
				}).Return(errors.New("an error"))
				return &service{
					dbRepository:    mockDB,
					redisRepository: mockRedis,
					client:          mockClient,
				}
			},
		},
		{
			name: "given happy flow should return nil error",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
			mockFunc: func() *service {
				mockDB := &mocks.Repository{}
				mockRedis := &mocks.RedisRepository{}
				mockClient := &mocks2.Client{}

				mockClient.
					On("GetBerries", mock.Anything, api.BerriesRequest{}).
					Return(&api.BerriesResponse{
						Count:    0,
						Next:     "",
						Previous: "",
						Results: []api.Berry{
							{
								Name: "1",
								Url:  "1",
							},
						},
					}, nil)

				mockDB.On("CreateBerry", mock.Anything, []model.Berry{
					{
						Name: "1",
						URL:  "1",
					},
				}).Return(nil)
				return &service{
					dbRepository:    mockDB,
					redisRepository: mockRedis,
					client:          mockClient,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.mockFunc()
			if err := s.SyncData(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("SyncData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_service_GetItems(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		args     args
		want     *model.BerriesResponse
		wantErr  bool
		mockFunc func() *service
	}{
		{
			name: "given result from redis should return response and no error",
			args: args{
				ctx: context.Background(),
			},
			want: &model.BerriesResponse{
				Berries: []model.Berry{
					{
						Name: "1",
						URL:  "1",
					},
				},
			},
			wantErr: false,
			mockFunc: func() *service {
				mockDB := &mocks.Repository{}
				mockRedis := &mocks.RedisRepository{}
				mockClient := &mocks2.Client{}

				mockRedis.
					On("GetData", mock.Anything).
					Return(&model.BerriesResponse{Berries: []model.Berry{
						{
							Name: "1",
							URL:  "1",
						},
					}}, nil)
				return &service{
					dbRepository:    mockDB,
					redisRepository: mockRedis,
					client:          mockClient,
				}
			},
		},
		{
			name: "given nil result from redis but got an error when call FetchBerries" +
				" should return nil response and an error",
			args: args{
				ctx: context.Background(),
			},
			want:    nil,
			wantErr: true,
			mockFunc: func() *service {
				mockDB := &mocks.Repository{}
				mockRedis := &mocks.RedisRepository{}
				mockClient := &mocks2.Client{}

				mockRedis.
					On("GetData", mock.Anything).
					Return(nil, redis.Nil)

				mockDB.
					On("FetchBerries", mock.Anything).
					Return(nil, errors.New("an error"))

				return &service{
					dbRepository:    mockDB,
					redisRepository: mockRedis,
					client:          mockClient,
				}
			},
		},
		{
			name: "given success when call FetchBerries but got an error when SetData to redis" +
				" should return response and no error",
			args: args{
				ctx: context.Background(),
			},
			want: &model.BerriesResponse{
				Berries: []model.Berry{
					{
						Name: "1",
						URL:  "1",
					},
				},
			},
			wantErr: false,
			mockFunc: func() *service {
				mockDB := &mocks.Repository{}
				mockRedis := &mocks.RedisRepository{}
				mockClient := &mocks2.Client{}

				mockRedis.
					On("GetData", mock.Anything).
					Return(nil, redis.Nil)

				mockDB.
					On("FetchBerries", mock.Anything).
					Return(&model.BerriesResponse{
						Berries: []model.Berry{
							{
								Name: "1",
								URL:  "1",
							},
						},
					}, nil)

				mockRedis.On("SetData", mock.Anything, &model.BerriesResponse{
					Berries: []model.Berry{
						{
							Name: "1",
							URL:  "1",
						},
					},
				}).Return(errors.New("an error"))

				return &service{
					dbRepository:    mockDB,
					redisRepository: mockRedis,
					client:          mockClient,
				}
			},
		},
		{
			name: "given happy flow when GetItems" +
				" should return response and no error",
			args: args{
				ctx: context.Background(),
			},
			want: &model.BerriesResponse{
				Berries: []model.Berry{
					{
						Name: "1",
						URL:  "1",
					},
				},
			},
			wantErr: false,
			mockFunc: func() *service {
				mockDB := &mocks.Repository{}
				mockRedis := &mocks.RedisRepository{}
				mockClient := &mocks2.Client{}

				mockRedis.
					On("GetData", mock.Anything).
					Return(nil, redis.Nil)

				mockDB.
					On("FetchBerries", mock.Anything).
					Return(&model.BerriesResponse{
						Berries: []model.Berry{
							{
								Name: "1",
								URL:  "1",
							},
						},
					}, nil)

				mockRedis.On("SetData", mock.Anything, &model.BerriesResponse{
					Berries: []model.Berry{
						{
							Name: "1",
							URL:  "1",
						},
					},
				}).Return(nil)

				return &service{
					dbRepository:    mockDB,
					redisRepository: mockRedis,
					client:          mockClient,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.mockFunc()
			got, err := s.GetItems(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetItems() got = %v, want %v", got, tt.want)
			}
		})
	}
}
