package api

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/inasknh/simple-poke-app/internal/config"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_client_GetBerries_Success(t *testing.T) {
	r := resty.New()

	// activate mock
	httpmock.ActivateNonDefault(r.GetClient())
	defer httpmock.DeactivateAndReset()

	respSuccess := &BerriesResponse{
		Count:    2,
		Next:     "",
		Previous: "",
		Results: []Berry{
			{
				Name: "leppa",
				Url:  "https://pokeapi.co/api/v2/berry/6/",
			},
			{
				Name: "oran",
				Url:  "https://pokeapi.co/api/v2/berry/7/",
			},
		},
	}
	// mock response
	httpmock.RegisterResponder("GET",
		"https://pokeapi.co/api/v2/berry?limit=10&offset=0",
		httpmock.NewJsonResponderOrPanic(http.StatusOK, respSuccess))

	// Inject Resty into client
	c := NewClient(config.Api{
		Host:     "https://pokeapi.co/api/v2/",
		BerryURL: "berry",
	}, r)

	resp, err := c.GetBerries(context.Background(), BerriesRequest{
		Limit:  10,
		Offset: 0,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, resp.Count)
}

func Test_client_GetBerries_ServerErrorThenRetry(t *testing.T) {
	// Create Resty with retry config
	r := resty.New().
		SetRetryCount(2).AddRetryCondition(func(response *resty.Response, err error) bool {
		return err != nil || response.StatusCode() >= 500
	})

	httpmock.ActivateNonDefault(r.GetClient())
	defer httpmock.DeactivateAndReset()

	respSuccess := &BerriesResponse{
		Count:    1,
		Next:     "",
		Previous: "",
		Results: []Berry{
			{
				Name: "leppa",
				Url:  "https://pokeapi.co/api/v2/berry/6/",
			},
		},
	}

	// First responder: 500, second responder: success
	callCount := 0
	httpmock.RegisterResponder("GET",
		"https://pokeapi.co/api/v2/berry?limit=10&offset=0",
		func(req *http.Request) (*http.Response, error) {
			callCount++
			if callCount == 1 {
				return httpmock.NewStringResponse(500, "internal error"), nil
			}
			return httpmock.NewJsonResponse(200, respSuccess)
		},
	)

	// Inject Resty into client
	c := NewClient(config.Api{
		Host:     "https://pokeapi.co/api/v2/",
		BerryURL: "berry",
	}, r)

	resp, err := c.GetBerries(context.Background(), BerriesRequest{
		Limit:  10,
		Offset: 0,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, resp.Count)
	assert.Equal(t, 2, callCount) // retried once
}
