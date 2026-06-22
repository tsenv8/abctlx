package sources

import (
	"abctlx/helpers"
	"abctlx/internal/airbyte"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Service interface {
	CreateSource(ctx context.Context, params CreateSourceParams, workspaceId string, token string) *CreateSourceResponse
	UpdateSource(ctx context.Context, params *UpdateSourceRequest, sourceName string, token string) *UpdateSourceResponse
	DeleteSource(ctx context.Context, sourceName string, token string) bool
	ListSources(ctx context.Context, token string) *ListSourcesResponse
	GetSourceId(ctx context.Context, name string, token string) (*SourceData, error)
}

type service struct {
	client airbyte.AirbyteClient
}

func NewService(c airbyte.AirbyteClient) *service {
	return &service{
		client: c,
	}
}

func (s *service) CreateSource(
	ctx context.Context,
	params CreateSourceParams,
	workspaceId string,
	token string,
) *CreateSourceResponse {
	var response CreateSourceResponse

	cdcReplicationMethod := CDCReplicationMethodParameter{
		Method:          "CDC",
		Plugin:          "pgoutput",
		ReplicationSlot: params.ReplicationSlot,
		Publication:     params.PublicationName,
	}

	sourcePostgresConf := PostgresConfigurationParameter{
		SourceType: "postgres",
		Host:       params.HostName,
		Port:       params.Port,
		Database:   params.DBName,
		Schemas:    params.Schemas,
		Username:   params.Username,
		Password:   params.Password,
		SSlMode: &SSLModeParameter{
			Mode: "disable",
		},
		ReplicationMethod: cdcReplicationMethod,
		TunnelMethod: TunnelMethodParameter{
			TunnelMethod: "NO_TUNNEL",
		},
	}

	requestBody := CreateSourceRequest{
		Name:          params.Name,
		WorkspaceId:   workspaceId,
		Configuration: sourcePostgresConf,
	}

	res, err := s.client.Request(
		ctx,
		http.MethodPost,
		airbyte.SOURCES_ENDPOINT,
		requestBody,
		&token,
	)

	if err != nil {
		helpers.Error("Create Source Request Failed", err)
	}

	err = json.Unmarshal(res.Body, &response)
	if err != nil {
		helpers.Error("Create Source JSON Unmarshal Failed", err)
	}

	return &response
}

func (s *service) UpdateSource(
	ctx context.Context,
	params *UpdateSourceRequest,
	sourceName string,
	token string,
) *UpdateSourceResponse {
	response := UpdateSourceResponse{}
	// token := s.GetAccessToken()
	source, err := s.GetSourceId(ctx, sourceName, token)

	if err != nil {
		helpers.Error("Source Not Found", err)
	}

	if source.SourceId == "" {
		helpers.Error("Source Not Found", err)
	}

	// pretty.Print(params)
	req, err := s.client.Request(
		ctx,
		http.MethodPatch,
		airbyte.SOURCES_ENDPOINT+"/"+source.SourceId,
		params,
		&token,
	)

	if err != nil {
		helpers.Error("Update Source Request Failed", err)
	}

	err = json.Unmarshal(req.Body, &response)
	if err != nil {
		helpers.Error("Update Source JSON Unmarshal Failed", err)
	}

	return &response
}

func (s *service) DeleteSource(ctx context.Context, sourceName string, token string) bool {
	// token := s.GetAccessToken()
	source, err := s.GetSourceId(ctx, sourceName, token)
	if err != nil {
		helpers.Error("Get Source Id Request Failed", err)
	}

	req, err := s.client.Request(
		ctx,
		http.MethodDelete,
		airbyte.SOURCES_ENDPOINT+"/"+*&source.SourceId,
		nil,
		&token,
	)

	if err != nil {
		helpers.Error("Delete Source Request Failed", err)
	}

	if req.Status >= 400 {
		return false
	}

	return true
}

func (s *service) ListSources(ctx context.Context, token string) *ListSourcesResponse {
	var response ListSourcesResponse
	// token := s.GetAccessToken()

	req, err := s.client.Request(
		ctx,
		http.MethodGet,
		airbyte.SOURCES_ENDPOINT,
		nil,
		&token,
	)

	if err != nil {
		helpers.Error("List Sources Request Failed", err)
	}

	if req == nil {
		helpers.Error("List Sources Request Failed", err)
	}

	err = json.Unmarshal(req.Body, &response)
	if err != nil {
		helpers.Error("List Sources Unmarshal Failed", err)
	}

	return &response
}

func (s *service) GetSourceId(ctx context.Context, name string, token string) (*SourceData, error) {
	sources := s.ListSources(ctx, token)
	var targetSource SourceData
	var sourceId *string

	for _, source := range sources.Data {
		if source.Name == name {
			targetSource = source
			break
		}
	}

	if sourceId == nil {
		return nil, fmt.Errorf("Source ID not found.")
	}

	// pretty.Print(targetSource)

	return &targetSource, nil
}
