package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type freteRapidoService struct {
	BaseURL    string
	HttpClient http.Client
}

type FreteRapidoQuoteSimulationRequest struct {
	Recipient struct {
		Zipcode int    `json:"zipcode"`
		Country string `json:"country"`
		Type    int    `json:"type"`
	} `json:"recipient"`

	Shipper struct {
		RegisteredNumber string `json:"registered_number"`
		PlatformCode     string `json:"platform_code"`
		Token            string `json:"token"`
	} `json:"shipper"`
	Dispatchers []struct {
		RegisteredNumber string `json:"registered_number"`
		Zipcode          int    `json:"zipcode"`
		Volumes          []struct {
			Category      string  `json:"category"`
			Amount        int     `json:"amount"`
			UnitaryWeight float64 `json:"unitary_weight"`
			UnitaryPrice  float64 `json:"unitary_price"`
			Sku           string  `json:"sku"`
			Height        float64 `json:"height"`
			Width         float64 `json:"width"`
			Length        float64 `json:"length"`
		} `json:"volumes"`
	} `json:"dispatchers"`
	SimulationType []int `json:"simulation_type"`
}

type FreteRapidoQuoteSimulationResponse struct {
	Carrier []struct {
		Name     string  `json:"name"`
		Service  string  `json:"service"`
		Deadline int     `json:"deadline"`
		Price    float64 `json:"price"`
	} `json:"carrier"`
}

func NewFreteRapidoService(url string) *freteRapidoService {
	var client = http.Client{}

	return &freteRapidoService{
		url,
		client,
	}
}

func (s *freteRapidoService) Quote(qr *FreteRapidoQuoteSimulationRequest) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/quote/simulate", s.BaseURL)

	jsonBody, err := json.Marshal(qr)

	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}
	defer request.Body.Close()

	resp, err := s.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}

	bodyResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	keyVal := make(map[string]interface{})
	json.Unmarshal(bodyResponse, &keyVal)

	return keyVal, nil
}
