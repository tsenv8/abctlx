package airbyte

import (
	"abctlx/internal/config"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kr/pretty"
)

type AirbyteService interface {
	// General
	GenerateAccessToken() *GenerateAccessTokenResponse
	SetClientToken()
	Health() *HealthCheckResponse

	//Sources
	CreateSource() *CreateSourceResponse
	ListSources() (*AbctlxResponse, error)
	GetSourceId()

	//Workspace
	ListWorkspaces() *ListWorkspacesResponse
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

func (s *airbyteService) GetWorkspaceId() *string {
	return &s.ListWorkspaces().Data[0].WorkspaceId
}

func (s *airbyteService) ListWorkspaces() *ListWorkspacesResponse {
	var response ListWorkspacesResponse
	token := s.GetAccessToken()
	res, err := s.client.Request(
		s.ctx,
		http.MethodGet,
		LIST_WORKSPACES_ENDPOINT,
		nil,
		&token,
	)
	if err != nil {
		NewAirbyteError(REQUEST_FAIL, "List Workspaces", err).Print()
	}

	err = json.Unmarshal(res.Body, &response)
	if err != nil {
		NewAirbyteError(JSON_UNMARSHAL_FAIL, "List Workspaces", err).Print()
	}

	return &response
}

func (s *airbyteService) CreateSource(params CreateSourceParams) *CreateSourceResponse {
	var response CreateSourceResponse
	workspaceId := s.GetWorkspaceId()
	token := s.GetAccessToken()

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
		WorkspaceId:   *workspaceId,
		Configuration: sourcePostgresConf,
	}

	res, err := s.client.Request(
		s.ctx,
		http.MethodPost,
		CREATE_SOURCE_ENDPOINT,
		requestBody,
		&token,
	)

	if err != nil {
		NewAirbyteError(REQUEST_FAIL, "Create Source", err).Print()
	}

	err = json.Unmarshal(res.Body, &response)
	if err != nil {
		NewAirbyteError(JSON_UNMARSHAL_FAIL, "Create Source", err).Print()
	}

	return &response
}

func (s *airbyteService) UpdateSource(params *UpdateSourceRequest, targetsource string) *UpdateSourceResponse {
	response := UpdateSourceResponse{}
	token := s.GetAccessToken()
	sourceId, err := s.GetSourceId(targetsource)

	if err != nil {
		NewAirbyteError("No such source found.", "Source Id", err).Print()
	}

	if sourceId == nil {
		NewAirbyteError("No such source found.", "Source Id", nil).Print()
	}

	pretty.Print(params)
	req, err := s.client.Request(
		s.ctx,
		http.MethodPatch,
		SOURCES_ENDPOINT+"/"+*sourceId,
		params,
		&token,
	)

	if err != nil {
		NewAirbyteError(REQUEST_FAIL, "Update Source", err).Print()
	}

	err = json.Unmarshal(req.Body, &response)
	if err != nil {
		NewAirbyteError(JSON_UNMARSHAL_FAIL, "Update Source", err).Print()
	}

	return &response
}

func (s *airbyteService) DeleteSource(sourceName string) bool {
	token := s.GetAccessToken()
	sourceId, err := s.GetSourceId(sourceName)
	if err != nil {
		NewAirbyteError(REQUEST_FAIL, "Source Id", err).Print()
	}

	req, err := s.client.Request(
		s.ctx,
		http.MethodDelete,
		SOURCES_ENDPOINT+"/"+*sourceId,
		nil,
		&token,
	)

	if err != nil {
		NewAirbyteError(REQUEST_FAIL, "Delete Source", err).Print()
	}

	if req.Status >= 400 {
		return false
	}

	return true
}

func (s *airbyteService) ListSources() *ListSourcesResponse {
	var response ListSourcesResponse
	token := s.GetAccessToken()

	req, err := s.client.Request(
		s.ctx,
		http.MethodGet,
		LIST_SOURCES_ENDPOINT,
		nil,
		&token,
	)

	if err != nil {
		NewAirbyteError(REQUEST_FAIL, "List Sources", err).Print()
	}

	if req == nil {
		NewAirbyteError(REQUEST_FAIL, "List Sources", err).Print()
	}

	err = json.Unmarshal(req.Body, &response)
	if err != nil {
		NewAirbyteError(JSON_UNMARSHAL_FAIL, "List Sources", err)
	}

	return &response
}

func (s *airbyteService) Health() *HealthCheckResponse {
	_, err := s.client.Request(
		s.ctx,
		http.MethodGet,
		HEALTH_CHECK_ENDPOINT,
		nil,
		nil,
	)

	if err != nil {
		return &HealthCheckResponse{
			Status: false,
		}
	}

	return &HealthCheckResponse{
		Status: true,
	}
}

func (s *airbyteService) GenerateAccessToken() *GenerateAccessTokenResponse {
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
		nil,
	)

	if err != nil {
		NewAirbyteError(REQUEST_FAIL, "Generate Access Token", err).Print()
	}

	err = json.Unmarshal(res.Body, &response)
	if err != nil {
		NewAirbyteError(JSON_UNMARSHAL_FAIL, "Generate Access Token", err).Print()
	}

	return &response
}

func (s *airbyteService) GetAccessToken() string {
	res := s.GenerateAccessToken()
	s.client.SetToken(res.AccessToken)
	return s.client.GetToken()
}

func (s *airbyteService) GetSourceId(name string) (*string, error) {
	sources := s.ListSources()
	var targetSource ListSourcesData
	var sourceId *string

	for _, source := range sources.Data {
		if source.Name == name {
			targetSource = source
			sourceId = &source.SourceId
			break
		}
	}

	if sourceId == nil {
		return nil, fmt.Errorf("Source ID not found.")
	}

	pretty.Print(name)
	pretty.Print(targetSource)

	return sourceId, nil
}
