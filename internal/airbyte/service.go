package airbyte

import (
	"abctlx/internal/config"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type AirbyteService interface {
	// General
	GenerateAccessToken() (*GenerateAccessTokenResponse, error)
	Health() (*HealthCheckResponse, error)

	//Sources
	// ListSources() (*AbctlxResponse, error)

	//Workspace
	// ListWorkspaces() (*AbctlxResponse, error)
}

type airbyteService struct {
	ctx    context.Context
	client AirbyteClient
}

func NewAirbyteService(ctx context.Context) *airbyteService {
	abClient := NewAirbyteClient(config.Data)
	return &airbyteService{
		ctx:    ctx,
		client: abClient,
	}
}

func (s *airbyteService) Health() (*HealthCheckResponse, error) {

	res, err := s.client.Request(
		s.ctx,
		http.MethodGet,
		HEALTH_CHECK_ENDPOINT,
		nil,
	)

	if err != nil {
		return &HealthCheckResponse{
			Status: false,
		}, err
	}

	log.Println(res)
	return &HealthCheckResponse{
		Status: true,
	}, nil
}

func (s *airbyteService) GenerateAccessToken() (*GenerateAccessTokenResponse, error) {
	var response GenerateAccessTokenResponse
	cfg := s.client.GetConfig()
	tokenRequest := GenerateAccessTokenRequest{
		ClientId:  cfg.ClientId,
		ClientKey: cfg.ClientKey,
	}

	res, err := s.client.Request(
		s.ctx,
		http.MethodPost,
		GENERATE_ACCESS_TOKEN_ENDPOINT,
		tokenRequest,
	)

	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(res.Body, &response)
	if err != nil {
		return nil, fmt.Errorf(JSON_UNMARSHAL_FAIL)
	}

	return &response, nil
}
