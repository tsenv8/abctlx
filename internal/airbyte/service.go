package airbyte

import (
	"abctlx/internal/config"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kr/pretty"
)

type AirbyteService interface {
	// General
	generateAccessToken() *GenerateAccessTokenResponse
	GetAccessToken() string
	Health() *HealthCheckResponse

	//Sources
	CreateSource(params CreateSourceParams) *CreateSourceResponse
	UpdateSource(params *UpdateSourceRequest, sourceName string) *UpdateSourceResponse
	DeleteSource(sourceName string) bool
	ListSources() *ListSourcesResponse
	GetSourceId(name string) (*SourceData, error)

	//Destinations
	CreateDestination(params CreateDestinationRequest) DestinationData
	UpdateDestination(params UpdateDestinationRequest) DestinationData
	DeleteDestination(destName string) bool
	ListDestinations(limit int) ListDestinationResponse
	GetDestination(destName string) DestinationData

	//Connections
	CreateConnection(params CreateConnectionRequest) ConnectionData
	UpdateConnection(params UpdateConnectionRequest, connectionName string) ConnectionData
	DeleteConnection(connectionName string) bool
	ListConnections(limit *int) ListConnectionResponse
	GetConnection(connectionName string) ConnectionData

	//Workspace
	ListWorkspaces() *ListWorkspacesResponse
	GetWorkspaceId() *string
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

// Connections
func (s *airbyteService) GetConnection(connectionName string) ConnectionData {
	connections := s.ListConnections(nil)
	var targetConnection ConnectionData
	for _, connection := range connections.Data {
		if connection.Name == connectionName {
			targetConnection = connection
			break
		}
	}

	pretty.Print(targetConnection)
	return targetConnection
}

func (s *airbyteService) ListConnections(limit *int) ListConnectionResponse {
	var response ListConnectionResponse
	var endpoint string
	token := s.GetAccessToken()

	if limit != nil {
		endpoint = CONNECTION_ENDPOINT + "?limit=" + strconv.Itoa(*limit)
	} else {
		endpoint = CONNECTION_ENDPOINT
	}

	res, err := s.client.Request(
		s.ctx,
		http.MethodGet,
		endpoint,
		nil,
		&token,
	)
	if err != nil {
		NewAirbyteError(REQUEST_FAIL, "List Connections", err)
	}

	err = json.Unmarshal(res.Body, &response)
	if err != nil {
		NewAirbyteError(JSON_UNMARSHAL_FAIL, "List Connections", err)
	}

	return response
}

func (s *airbyteService) DeleteConnection(connectionName string) bool {
	token := s.GetAccessToken()
	req, err := s.client.Request(
		s.ctx,
		http.MethodDelete,
		CONNECTION_ENDPOINT,
		nil,
		&token,
	)
	if err != nil {
		NewAirbyteError(REQUEST_FAIL, "Delete Connection", err).Print()
	}

	if req.Status >= http.StatusBadRequest {
		return false
	}

	return true
}

func (s *airbyteService) UpdateConnection(params UpdateConnectionRequest, connectionName string) ConnectionData {
	var response ConnectionData
	token := s.GetAccessToken()
	req, err := s.client.Request(
		s.ctx,
		http.MethodPatch,
		CONNECTION_ENDPOINT+"/"+connectionName,
		params,
		&token,
	)

	if err != nil {
		NewAirbyteError(REQUEST_FAIL, "Update Connection", err).Print()
	}

	err = json.Unmarshal(req.Body, &response)
	if err != nil {
		NewAirbyteError(JSON_UNMARSHAL_FAIL, "Update Connection", err).Print()
	}

	return response
}

func (s *airbyteService) CreateConnection(params CreateConnectionRequest) ConnectionData {
	var response ConnectionData
	token := s.GetAccessToken()
	req, err := s.client.Request(
		s.ctx,
		http.MethodPost,
		CONNECTION_ENDPOINT,
		params,
		&token,
	)

	if err != nil {
		NewAirbyteError(REQUEST_FAIL, "Create Connection", err).Print()
	}

	err = json.Unmarshal(req.Body, &response)
	if err != nil {
		NewAirbyteError(JSON_UNMARSHAL_FAIL, "Create Connection", err).Print()
	}

	return response
}

// Destinations
func (s *airbyteService) UpdateDestination(flags UpdateDestinationFlags) DestinationData {
	var response DestinationData
	token := s.GetAccessToken()
	dest := s.GetDestination(*flags.DestName)

	config := DestinationConfigurationParameter{
		Host:     *flags.Host,
		Port:     *flags.Port,
		Database: *flags.Database,
		Username: *flags.Username,
		Password: *flags.Password,
	}

	updateDestReq := UpdateDestinationRequest{
		Name:          *flags.Name,
		Configuration: config,
	}

	req, err := s.client.Request(
		s.ctx,
		http.MethodPatch,
		DESTINATION_ENDPOINT+"/"+dest.DestinationId,
		updateDestReq,
		&token,
	)

	if err != nil {
		NewAirbyteError(REQUEST_FAIL, "Update Destination", err).Print()
	}

	err = json.Unmarshal(req.Body, &response)
	if err != nil {
		NewAirbyteError(JSON_UNMARSHAL_FAIL, "Update Destination", err).Print()
	}

	return response
}

func (s *airbyteService) DeleteDestination(destName string) bool {
	destination := s.GetDestination(destName)
	token := s.GetAccessToken()
	res, err := s.client.Request(
		s.ctx,
		http.MethodDelete,
		DESTINATION_ENDPOINT+"/"+destination.DestinationId,
		nil,
		&token,
	)
	if err != nil {
		NewAirbyteError(REQUEST_FAIL, "Delete Destination", err)
	}

	if res.Status >= http.StatusBadRequest {
		return false
	}

	return true
}

func (s *airbyteService) CreateDestination(flags CreateDestinationFlags) DestinationData {
	var response DestinationData
	token := s.GetAccessToken()
	workspaceId := s.GetWorkspaceId()

	config := DestinationConfigurationParameter{
		Host:     "localhost",
		Port:     8123,
		Database: "chdb",
		Username: "default",
		Password: "1",
		TunnelMethod: TunnelMethodParameter{
			TunnelMethod: "NO_TUNNEL",
		},
		DestinationType: "clickhouse",
	}

	createDestReq := CreateDestinationRequest{
		Name:          flags.Name,
		WorkspaceId:   *workspaceId,
		Configuration: config,
	}

	req, err := s.client.Request(
		s.ctx,
		http.MethodPost,
		DESTINATION_ENDPOINT,
		createDestReq,
		&token,
	)

	if err != nil {
		NewAirbyteError(REQUEST_FAIL, "Create Destination", err).Print()
	}

	err = json.Unmarshal(req.Body, &response)
	if err != nil {
		NewAirbyteError(JSON_UNMARSHAL_FAIL, "Create Destination", err).Print()
	}

	return response
}

func (s *airbyteService) GetDestination(destName string) DestinationData {
	destinations := s.ListDestinations(nil)
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

func (s *airbyteService) ListDestinations(limit *int) ListDestinationResponse {
	var response ListDestinationResponse
	var finalEndpoint string
	token := s.GetAccessToken()

	if limit != nil {
		finalEndpoint = DESTINATION_ENDPOINT + "?limit=" + strconv.Itoa(*limit)
	} else {
		finalEndpoint = DESTINATION_ENDPOINT
	}

	req, err := s.client.Request(
		s.ctx,
		http.MethodGet,
		finalEndpoint,
		nil,
		&token,
	)
	if err != nil {
		NewAirbyteError(REQUEST_FAIL, "List Destinations", err).Print()
	}

	err = json.Unmarshal(req.Body, &response)
	if err != nil {
		NewAirbyteError(JSON_UNMARSHAL_FAIL, "List Destinations", err).Print()
	}

	return response
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
		SOURCES_ENDPOINT,
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

func (s *airbyteService) UpdateSource(params *UpdateSourceRequest, sourceName string) *UpdateSourceResponse {
	response := UpdateSourceResponse{}
	token := s.GetAccessToken()
	source, err := s.GetSourceId(sourceName)

	if err != nil {
		NewAirbyteError("No such source found.", "Source Id", err).Print()
	}

	if source.SourceId == "" {
		NewAirbyteError("No such source found.", "Source Id", nil).Print()
	}

	pretty.Print(params)
	req, err := s.client.Request(
		s.ctx,
		http.MethodPatch,
		SOURCES_ENDPOINT+"/"+*&source.SourceId,
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
	source, err := s.GetSourceId(sourceName)
	if err != nil {
		NewAirbyteError(REQUEST_FAIL, "Source Id", err).Print()
	}

	req, err := s.client.Request(
		s.ctx,
		http.MethodDelete,
		SOURCES_ENDPOINT+"/"+*&source.SourceId,
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
		SOURCES_ENDPOINT,
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

func (s *airbyteService) generateAccessToken() *GenerateAccessTokenResponse {
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
	res := s.generateAccessToken()
	s.client.SetToken(res.AccessToken)
	return s.client.GetToken()
}

func (s *airbyteService) GetSourceId(name string) (*SourceData, error) {
	sources := s.ListSources()
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

	pretty.Print(targetSource)

	return &targetSource, nil
}
