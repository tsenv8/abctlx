package airbyte

import (
	"abctlx/internal/config"
	"fmt"
	"io"
	"log"
	"net/http"
)

type AirbyteRepository interface {
	GenerateAccessToken() string
	HealthCheck() string
}

type airbyteRepository struct {
	config config.AirbyteConfig
}

func New(c config.AirbyteConfig) AirbyteRepository {
	return &airbyteRepository{config: c}
}

func (r *airbyteRepository) GenerateAccessToken() string {
	url := r.config.URL + "/v1/applications/token"
	req, _ := http.NewRequest("POST", url, nil)

	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("Failed Command")
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	return string(body)
}

func (r *airbyteRepository) HealthCheck() string {

	url := r.config.URL + "/v1/health"
	req, _ := http.NewRequest("GET", url, nil)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("Failed Command")
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))
	return string(body)
}
