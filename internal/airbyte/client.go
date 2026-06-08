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

	"github.com/kr/pretty"
)

type AirbyteClient interface {
	// General
	// GenerateAccessToken() (*AbctlxResponse, error)
	// HealthCheck() (*AbctlxResponse, error)
	Request(ctx context.Context, method string, endpoint string, requestBody any, token *string) (*AbctlxResponse, error)
	GetURL(endpoint *string) string
	GetConfig() config.AirbyteConfig
	SetToken(string) string
	GetToken() string

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

func (c *airbyteClient) SetToken(token string) string {
	c.AccessToken = &token
	return *c.AccessToken
}

func (c *airbyteClient) GetToken() string {
	return *c.AccessToken
}

func (ac *airbyteClient) Request(
	ctx context.Context,
	method string,
	endpoint string,
	requestBody any,
	token *string,
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

	pretty.Print("\n Making a request with: [" + method + "]" + " " + url)

	req, err := http.NewRequestWithContext(ctx, method, url, buf)
	if err != nil {
		return nil, err
	}

	//headers and tokens
	if buf != nil {
		req.Header.Set("Content-Type", HEADER_CONTENT_TYPE)
	}

	req.Header.Set("Accept", HEADER_CONTENT_TYPE)
	if token != nil {
		req.Header.Set("Authorization", "Bearer "+*token)
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
		Status:   res.StatusCode,
	}, nil
}
