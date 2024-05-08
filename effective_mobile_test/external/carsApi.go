package external

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type CarsAPI struct {
	baseUrl string
}

type Owner struct {
	Name       string
	Surname    string
	Patronymic string
}

type CarDTO struct {
	RegNum string
	Mark   string
	Model  string
	Year   int
	Owner  Owner
}

func NewCarsAPI() *CarsAPI {
	carsExternalApiURL := os.Getenv("EXTERNAL_CAR_API_URL")
	if carsExternalApiURL == "" {
		log.Fatal("External API url for cars not provided in .env file")
	}
	return &CarsAPI{baseUrl: carsExternalApiURL}
}

func (api *CarsAPI) Get(subdirectory, queryString string) ([]CarDTO, error) {
	fullUrl := api.baseUrl + "/" + subdirectory + "?" + queryString
	res, err := http.Get(fullUrl)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result []CarDTO

	if err = json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}
