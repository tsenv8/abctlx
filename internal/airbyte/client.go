package airbyte

import (
	"abctlx/internal/config"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type AirbyteClient interface {
	// General
	// GenerateAccessToken() (*AbctlxResponse, error)
	// HealthCheck() (*AbctlxResponse, error)
	Request(ctx context.Context, method string, endpoint string, requestBody any) (*AbctlxResponse, error)
	GetURL(endpoint *string) string
	GetConfig() config.AirbyteConfig

	//Sources
	// ListSources() (*AbctlxResponse, error)

	//Workspace
	// ListWorkspaces() (*AbctlxResponse, error)
}

type airbyteClient struct {
	Config      config.AirbyteConfig
	Http        *http.Client
	AccessToken *string
}

func NewAirbyteClient(c config.AirbyteConfig) AirbyteClient {
	return &airbyteClient{
		Config: c,
		Http: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

func (ac *airbyteClient) GetURL(endpoint *string) string {
	port := ":" + strconv.Itoa(ac.Config.Port)

	if endpoint != nil {
		return ac.Config.URL + port + ac.Config.API_URL + *endpoint
	}

	return ac.Config.URL + port + ac.Config.API_URL
}

func (ac *airbyteClient) GetConfig() config.AirbyteConfig {
	return ac.Config
}

func (ac *airbyteClient) Request(
	ctx context.Context,
	method string,
	endpoint string,
	requestBody any,
) (*AbctlxResponse, error) {
	var buf io.Reader

	//marshal struct
	if requestBody != nil {
		jsonBytes, err := json.Marshal(requestBody)
		if err != nil {
			return nil, err
		}

		buf = bytes.NewBuffer(jsonBytes)
	}

	//build url
	url := ac.GetURL(&endpoint)
	req, err := http.NewRequestWithContext(ctx, method, url, buf)
	if err != nil {
		return nil, err
	}

	//headers and tokens
	if buf != nil {
		req.Header.Set("Content-Type", HEADER_CONTENT_TYPE)
	}
	
	req.Header.Set("Accept", HEADER_CONTENT_TYPE)
	if ac.AccessToken != nil {
		req.Header.Set("Authorization", "Bearer "+*ac.AccessToken)
	}

	res, err := ac.Http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		errBody, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("API Error: %s", string(errBody))
	}

	//res.Body reading
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &AbctlxResponse{
		Msg:      REQUEST_SUCCESS,
		Body:     body,
		Endpoint: &url,
	}, nil
}

// func (c *airbyteClient) GenerateAccessToken() (*AbctlxResponse, error) {
// 	var gatResponse GenerateAccessTokenResponse
// 	endpoint := GENERATE_ACCESS_TOKEN_ENDPOINT
// 	url := c.GetURL(&endpoint)

// 	tokenRequest := GenerateAccessTokenRequest{
// 		ClientId:  c.Config.ClientId,
// 		ClientKey: c.Config.ClientKey,
// 	}

// 	payload, err := json.Marshal(tokenRequest)
// 	if err != nil {
// 		return nil, err
// 	}

// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
// 	req.Header.Add("content-type", HEADER_CONTENT_TYPE)
// 	req.Header.Add("accept", HEADER_CONTENT_TYPE)

// 	if err != nil {
// 		return nil, fmt.Errorf(REQUEST_FAIL)
// 	}

// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf(REQUEST_FAIL)
// 	}

// 	defer res.Body.Close()
// 	body, _ := io.ReadAll(res.Body)

// 	err = json.Unmarshal(body, &gatResponse)
// 	if err != nil {
// 		return nil, fmt.Errorf(REQUEST_FAIL)
// 	}

// 	response := &AbctlxResponse{
// 		Msg:      REQUEST_SUCCESS,
// 		Body:     body,
// 		Endpoint: &url,
// 	}

// 	return response, nil
// }

// func (c *airbyteClient) HealthCheck() (*AbctlxResponse, error) {
// 	endpoint := HEALTH_CHECK_ENDPOINT
// 	url := c.GetURL(&endpoint)
// 	req, _ := http.NewRequest("GET", url, nil)

// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf(REQUEST_FAIL)
// 	}

// 	defer res.Body.Close()
// 	body, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf(REQUEST_FAIL)
// 	}

// 	return &AbctlxResponse{
// 		Msg:      REQUEST_SUCCESS,
// 		Body:     body,
// 		Endpoint: &url,
// 	}, nil
// }

// func (c *airbyteClient) ListSources() (*AbctlxResponse, error) {
// 	endpoint := LIST_SOURCES_ENDPOINT
// 	url := c.GetURL(&endpoint)
// 	req, _ := http.NewRequest("GET", url, nil)

// 	req.Header.Add("accept", HEADER_CONTENT_TYPE)

// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf(REQUEST_FAIL)
// 	}

// 	defer res.Body.Close()
// 	body, _ := io.ReadAll(res.Body)

// 	fmt.Println(string(body))
// 	return &AbctlxResponse{
// 		Msg:      REQUEST_SUCCESS,
// 		Body:     body,
// 		Endpoint: &url,
// 	}, nil
// }

// func (c *airbyteClient) CreateSource(params CreateSourceParams) (*AbctlxResponse, error) {
// 	endpoint := CREATE_SOURCE_ENDPOINT
// 	url := c.GetURL(&endpoint)

// 	var workspace WorkspaceData
// 	workspaceData, err := c.ListWorkspaces()
// 	if err != nil {
// 		return nil, fmt.Errorf(REQUEST_FAIL)
// 	}

// 	json.Unmarshal(workspaceData.Body, &workspace)

// 	cdcReplicationMethod := CDCReplicationMethodParameter{
// 		Method:          "CDC",
// 		Plugin:          "pgoutput",
// 		ReplicationSlot: params.replicationSlot,
// 		Publication:     params.publicationName,
// 	}

// 	sourcePostgresConf := PostgresConfigurationParameter{
// 		SourceType:        "postgres",
// 		Host:              params.hostName,
// 		Port:              params.port,
// 		Database:          params.dbName,
// 		Schemas:           params.schemas,
// 		Username:          params.username,
// 		Password:          params.password,
// 		SSlMode:           nil,
// 		ReplicationMethod: cdcReplicationMethod,
// 		TunnelMethod: TunnelMethodParameter{
// 			TunnelMethod: "NO_TUNNEL",
// 		},
// 	}

// 	sourceRequest := CreateSourceRequest{
// 		Name:          params.name,
// 		WorkspaceId:   workspace.WorkspaceId,
// 		Configuration: sourcePostgresConf,
// 	}

// 	payload, err := json.Marshal(&sourceRequest)
// 	if err != nil {
// 		return nil, fmt.Errorf(REQUEST_FAIL)
// 	}

// 	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))

// 	req.Header.Add("accept", HEADER_CONTENT_TYPE)
// 	req.Header.Add("content-type", HEADER_CONTENT_TYPE)

// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf(REQUEST_FAIL)
// 	}

// 	defer res.Body.Close()
// 	body, _ := io.ReadAll(res.Body)

// 	return &AbctlxResponse{
// 		Msg:      REQUEST_SUCCESS,
// 		Body:     body,
// 		Endpoint: &url,
// 	}, nil
// }

// func (c *airbyteClient) ListWorkspaces() (*AbctlxResponse, error) {
// 	var listWorkspacesResponse ListWorkspacesResponse
// 	url := LIST_WORKSPACES_ENDPOINT

// 	req, _ := http.NewRequest("GET", url, nil)
// 	req.Header.Add("accept", HEADER_CONTENT_TYPE)

// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf(REQUEST_FAIL)
// 	}

// 	defer res.Body.Close()
// 	body, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf(REQUEST_FAIL)
// 	}

// 	err = json.Unmarshal(body, &listWorkspacesResponse)
// 	if err != nil {
// 		return nil, fmt.Errorf(JSON_UNMARSHAL_FAIL)
// 	}

// 	fmt.Println(string(body))
// 	return &AbctlxResponse{
// 		Msg:      REQUEST_SUCCESS,
// 		Body:     body,
// 		Data:     listWorkspacesResponse,
// 		Endpoint: &url,
// 	}, nil
// }
