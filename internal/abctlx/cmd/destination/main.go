package destination

import (
	"abctlx/helpers"
	"abctlx/internal/airbyte"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kr/pretty"
)

type Service interface {
	CreateDestination(ctx context.Context, flags CreateDestinationFlags, token string, workspaceId string) DestinationData
	UpdateDestination(ctx context.Context, flags UpdateDestinationFlags, token string) DestinationData
	DeleteDestination(ctx context.Context, destName string, token string) bool
	ListDestinations(ctx context.Context, limit *int, token string) ListDestinationResponse
	GetDestination(ctx context.Context, destName string, token string) DestinationData
}

type service struct {
	client airbyte.AirbyteClient
}

func NewService(client airbyte.AirbyteClient) *service {
	return &service{
		client: client,
	}
}

func (s *service) UpdateDestination(ctx context.Context, flags UpdateDestinationFlags, token string) DestinationData {
	var response DestinationData
	// token := s.GetAccessToken()
	dest := s.GetDestination(ctx, flags.DestName, token)

	if dest.DestinationId == "" {
		helpers.Error("Update Destination Request Failed", nil)
		// NewAirbyteError(REQUEST_FAIL, "Update Destination", fmt.Errorf("No Destination Object Found")).Print()
	}

	// if flags.Host != "" {
	// 	config.Host = flags.Host
	// }

	// if flags.Port != "" {
	// 	config.Port = flags.Port
	// }

	// if flags.Database != "" {
	// 	config.Database = flags.Database
	// }

	// if flags.Username != "" {
	// 	config.Username = flags.Username
	// }

	// if flags.Username != "" {
	// 	config.Password = flags.Password
	// }

	config := DestinationConfigurationParameter{
		Host:     "",
		Port:     "",
		Database: "",
		Username: "",
		Password: "",
	}

	updateDestReq := UpdateDestinationRequest{
		Name:          flags.Name,
		Configuration: &config,
	}

	req, err := s.client.Request(
		ctx,
		http.MethodPatch,
		airbyte.DESTINATION_ENDPOINT+dest.DestinationId,
		updateDestReq,
		&token,
	)

	if err != nil {
		helpers.Error("Update Destination Request Failed", err)
	}

	err = json.Unmarshal(req.Body, &response)
	if err != nil {
		helpers.Error("Update Destination JSON Unmarshal Failed", err)
	}

	return response
}

func (s *service) DeleteDestination(ctx context.Context, destName string, token string) bool {
	destination := s.GetDestination(ctx, destName, token)
	// token := s.GetAccessToken()
	res, err := s.client.Request(
		ctx,
		http.MethodDelete,
		airbyte.DESTINATION_ENDPOINT+"/"+destination.DestinationId,
		nil,
		&token,
	)
	if err != nil {
		helpers.Error("Delete Destination Request Failed", err)
	}

	if res.Status >= http.StatusBadRequest {
		return false
	}

	return true
}

func (s *service) CreateDestination(ctx context.Context, flags CreateDestinationFlags, token string, workspaceId string) DestinationData {
	var response DestinationData
	var config DestinationConfigurationParameter
	// token := s.GetAccessToken()
	// workspaceId := s.GetWorkspaceId()

	if flags.ConfigType == "clickhouse" {
		config = DestinationConfigurationParameter{
			Host:     "localhost",
			Port:     "8123",
			Database: "chdb",
			Username: "default",
			Protocol: "http",
			Password: "1",
			TunnelMethod: TunnelMethodParameter{
				TunnelMethod: "NO_TUNNEL",
			},
			DestinationType: "clickhouse",
		}
	} else {
		helpers.Error("Create Destination Request Failed", fmt.Errorf("Failed to create config for %s", flags.Name))
	}

	createDestReq := CreateDestinationRequest{
		Name:          flags.Name,
		WorkspaceId:   workspaceId,
		Configuration: config,
	}

	req, err := s.client.Request(
		ctx,
		http.MethodPost,
		airbyte.DESTINATION_ENDPOINT,
		createDestReq,
		&token,
	)

	if err != nil {
		helpers.Error("Create Destination Request Failed", err)
	}

	err = json.Unmarshal(req.Body, &response)
	if err != nil {
		helpers.Error("Create Destination JSON Unmarshal Failed", err)
	}

	return response
}

func (s *service) GetDestination(ctx context.Context, destName string, token string) DestinationData {
	destinations := s.ListDestinations(ctx, nil, token)
	var targetDestination DestinationData
	for _, destination := range destinations.Data {
		if destination.Name == destName {
			targetDestination = destination
			break
		}
	}

	pretty.Print(targetDestination)
	return targetDestination
}

func (s *service) ListDestinations(ctx context.Context, limit *int, token string) ListDestinationResponse {
	var response ListDestinationResponse
	var finalEndpoint string
	// token := s.GetAccessToken()

	if limit != nil {
		finalEndpoint = airbyte.DESTINATION_ENDPOINT + "?limit=" + strconv.Itoa(*limit)
	} else {
		finalEndpoint = airbyte.DESTINATION_ENDPOINT
	}

	req, err := s.client.Request(
		ctx,
		http.MethodGet,
		finalEndpoint,
		nil,
		&token,
	)

	if err != nil {
		helpers.Error("List Destinations Request Failed", err)
		// NewAirbyteError(REQUEST_FAIL, "List Destinations", err).Print()
	}

	err = json.Unmarshal(req.Body, &response)
	if err != nil {
		helpers.Error("JSON Unmarshal Failed", err)
	}

	return response
}
