package airbyte

import (
	"abctlx/helpers"
	"abctlx/internal/config"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

// DryRun makes Request print every mutating call (anything but GET and HEAD)
// and skip sending it. The token request is exempt: it changes nothing and
// every other call needs it.
var DryRun bool

// ErrDryRun is returned in place of a response for a skipped call.
var ErrDryRun = errors.New("dry run: request not sent")

// dryRunExit ends the command after the first skipped call. A skipped call has
// no response, so nothing after it can run meaningfully. Tests replace it.
var dryRunExit = func() { os.Exit(0) }

func isMutation(method, endpoint string) bool {
	if endpoint == GENERATE_ACCESS_TOKEN_ENDPOINT {
		return false
	}
	return method != http.MethodGet && method != http.MethodHead
}

type AirbyteClient interface {
	Request(ctx context.Context, method string, endpoint string, requestBody any, token *string) (*AbctlxResponse, error)
	GetURL(endpoint *string) string
	GetConfig() config.AirbyteConfig
	SetToken(string) string
	GetToken() string
	// Table(header table.Row, body []table.Row)
	// ToRow(data []string) (table.Row, error)
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

	if DryRun && isMutation(method, endpoint) {
		fmt.Printf("\nDry run: not sent\n  %s %s\n", method, url)
		if requestBody != nil {
			if pretty, err := json.MarshalIndent(requestBody, "  ", "  "); err == nil {
				fmt.Printf("  %s\n", pretty)
			}
		}
		fmt.Println("\nNothing was changed.")
		dryRunExit()
		return nil, ErrDryRun
	}

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

	helpers.LogHttp(method, strconv.Itoa(res.StatusCode), url)
	// pretty.Printf("\n  %s  [ %s ] - %s ", color.HiGreenString("[ "+method+" ]"), strconv.Itoa(res.StatusCode), url)

	if res.StatusCode >= http.StatusBadRequest {
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
