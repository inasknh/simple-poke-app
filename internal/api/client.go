package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/inasknh/simple-poke-app/internal/config"
	"strconv"
)

type client struct {
	host       string
	berryURL   string
	rstyClient *resty.Client
}

type Client interface {
	GetBerries(ctx context.Context, request BerriesRequest) (*BerriesResponse, error)
}

func NewClient(config config.Api, rstyClient *resty.Client) Client {
	return &client{
		host:       config.Host,
		berryURL:   config.BerryURL,
		rstyClient: rstyClient,
	}
}

func (c *client) GetBerries(ctx context.Context, request BerriesRequest) (*BerriesResponse, error) {
	resp, err := c.rstyClient.
		//SetTimeout(5*time.Second).
		//SetRetryCount(3).
		//AddRetryCondition(func(response *resty.Response, err error) bool {
		//	return err != nil || response.StatusCode() >= 500 || response.StatusCode() == http.StatusTooManyRequests
		//}).
		R().
		SetContext(ctx).
		SetQueryParam("limit", strconv.Itoa(request.Limit)).
		SetQueryParam("offset", strconv.Itoa(request.Offset)).
		Get(fmt.Sprintf("%s%s", c.host, c.berryURL))

	if err != nil {
		return nil, err
	}

	var br BerriesResponse
	err = json.Unmarshal(resp.Body(), &br)
	if err != nil {
		return nil, err
	}
	return &br, nil

}
