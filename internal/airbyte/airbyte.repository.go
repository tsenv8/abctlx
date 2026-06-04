package airbyte

import (
	"abctlx/internal/config"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

type AirbyteRepository interface {
	GenerateAccessToken() string
	HealthCheck() string
	GetURL(*string) string
}

type airbyteRepository struct {
	config config.AirbyteConfig
}

func New(c config.AirbyteConfig) AirbyteRepository {
	return &airbyteRepository{config: c}
}

func (r *airbyteRepository) GetURL(addtlUrl *string) string {
	if addtlUrl != nil {
		return r.config.URL + strconv.Itoa(r.config.Port) + r.config.API_URL + *addtlUrl
	}
	return r.config.URL + strconv.Itoa(r.config.Port) + r.config.API_URL
}

func (r *airbyteRepository) GenerateAccessToken() string {
	apiUrl := "/applications/token"
	finalUrl := r.GetURL(&apiUrl)
	req, _ := http.NewRequest("POST", finalUrl, nil)

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

	url := r.config.URL + ":" + strconv.Itoa(r.config.Port) + "/v1/health"
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
