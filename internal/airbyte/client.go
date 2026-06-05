package airbyte

import (
	"abctlx/internal/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type AirbyteClient interface {
	// General
	GenerateAccessToken() (*AbctlxResponse, error)
	HealthCheck() (*AbctlxResponse, error)
	GetURL(*string) string

	//Sources
	ListSources() (*AbctlxResponse, error)

	//Workspace
	ListWorkspaces() (*AbctlxResponse, error)
}

type airbyteClient struct {
	config config.AirbyteConfig
}

func New(c config.AirbyteConfig) AirbyteClient {
	return &airbyteClient{config: c}
}

func (r *airbyteClient) GetURL(endpoint *string) string {
	port := ":" + strconv.Itoa(r.config.Port)

	if endpoint != nil {
		return r.config.URL + port + r.config.API_URL + *endpoint
	}
	return r.config.URL + port + r.config.API_URL
}

func (c *airbyteClient) GenerateAccessToken() (*AbctlxResponse, error) {
	endpoint := GENERATE_ACCESS_TOKEN_ENDPOINT
	url := c.GetURL(&endpoint)

	tokenRequest := GenerateAccessTokenRequest{
		ClientId:  c.config.ClientId,
		ClientKey: c.config.ClientKey,
	}

	payload, err := json.Marshal(tokenRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Add("content-type", HEADER_CONTENT_TYPE)
	req.Header.Add("accept", HEADER_CONTENT_TYPE)

	if err != nil {
		return nil, fmt.Errorf(REQUEST_FAIL)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(REQUEST_FAIL)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	return &AbctlxResponse{
		msg:      REQUEST_SUCCESS,
		body:     body,
		endpoint: &url,
	}, nil
}

func (c *airbyteClient) HealthCheck() (*AbctlxResponse, error) {
	endpoint := HEALTH_CHECK_ENDPOINT
	url := c.GetURL(&endpoint)
	req, _ := http.NewRequest("GET", url, nil)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(REQUEST_FAIL)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf(REQUEST_FAIL)
	}

	return &AbctlxResponse{
		msg:      REQUEST_SUCCESS,
		body:     body,
		endpoint: &url,
	}, nil
}

func (c *airbyteClient) ListSources() (*AbctlxResponse, error) {
	endpoint := LIST_SOURCES_ENDPOINT
	url := c.GetURL(&endpoint)
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", HEADER_CONTENT_TYPE)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(REQUEST_FAIL)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))
	return &AbctlxResponse{
		msg:      REQUEST_SUCCESS,
		body:     body,
		endpoint: &url,
	}, nil
}

func (c *airbyteClient) CreateSource(params CreateSourceParams) (*AbctlxResponse, error) {
	endpoint := CREATE_SOURCE_ENDPOINT
	url := c.GetURL(&endpoint)

	var workspace WorkspaceData
	workspaceData, err := c.ListWorkspaces()
	if err != nil {
		return nil, fmt.Errorf(REQUEST_FAIL)
	}

	json.Unmarshal(workspaceData.body, &workspace)

	cdcReplicationMethod := CDCReplicationMethodParameter{
		Method:          "CDC",
		Plugin:          "pgoutput",
		ReplicationSlot: params.replicationSlot,
		Publication:     params.publicationName,
	}

	sourcePostgresConf := PostgresConfigurationParameter{
		SourceType:        "postgres",
		Host:              params.hostName,
		Port:              params.port,
		Database:          params.dbName,
		Schemas:           params.schemas,
		Username:          params.username,
		Password:          params.password,
		SSlMode:           nil,
		ReplicationMethod: cdcReplicationMethod,
		TunnelMethod: TunnelMethodParameter{
			TunnelMethod: "NO_TUNNEL",
		},
	}

	sourceRequest := CreateSourceRequest{
		Name:          params.name,
		WorkspaceId:   workspace.WorkspaceId,
		Configuration: sourcePostgresConf,
	}

	payload, err := json.Marshal(&sourceRequest)
	if err != nil {
		return nil, fmt.Errorf(REQUEST_FAIL)
	}

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))

	req.Header.Add("accept", HEADER_CONTENT_TYPE)
	req.Header.Add("content-type", HEADER_CONTENT_TYPE)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(REQUEST_FAIL)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	return &AbctlxResponse{
		msg:      REQUEST_SUCCESS,
		body:     body,
		endpoint: &url,
	}, nil
}

func (c *airbyteClient) ListWorkspaces() (*AbctlxResponse, error) {
	var listWorkspacesResponse ListWorkspacesResponse
	url := LIST_WORKSPACES_ENDPOINT

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", HEADER_CONTENT_TYPE)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(REQUEST_FAIL)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf(REQUEST_FAIL)
	}

	err = json.Unmarshal(body, &listWorkspacesResponse)
	if err != nil {
		return nil, fmt.Errorf(JSON_UNMARSHAL_FAIL)
	}

	fmt.Println(string(body))
	return &AbctlxResponse{
		msg:      REQUEST_SUCCESS,
		body:     body,
		data:     listWorkspacesResponse,
		endpoint: &url,
	}, nil
}
