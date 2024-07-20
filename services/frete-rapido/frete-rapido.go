package freterapido

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Service struct {
	BaseURL    string
	HttpClient http.Client
}

type QuoteSimulationRequest struct {
	Recipient      Recipient `json:"recipient"`
	Shipper        Shipper   `json:"shipper"`
	Dispatchers    []Dispatcher
	SimulationType []int `json:"simulation_type"`
}

type Recipient struct {
	Zipcode int    `json:"zipcode"`
	Country string `json:"country"`
	Type    int    `json:"type"`
}

type Shipper struct {
	RegisteredNumber string `json:"registered_number"`
	PlatformCode     string `json:"platform_code"`
	Token            string `json:"token"`
}

type Dispatcher struct {
	RegisteredNumber string   `json:"registered_number"`
	Zipcode          int      `json:"zipcode"`
	Volumes          []Volume `json:"volumes"`
}

type Volume struct {
	Category      string  `json:"category"`
	Amount        int     `json:"amount"`
	UnitaryWeight float64 `json:"unitary_weight"`
	UnitaryPrice  float64 `json:"unitary_price"`
	Sku           string  `json:"sku"`
	Height        float64 `json:"height"`
	Width         float64 `json:"width"`
	Length        float64 `json:"length"`
}

type QuoteSimulationResponse struct {
	Dispatchers []struct {
		ID     string `json:"id"`
		Offers []struct {
			Carrier struct {
				CompanyName      string `json:"company_name"`
				Logo             string `json:"logo"`
				Name             string `json:"name"`
				Reference        int    `json:"reference"`
				RegisteredNumber string `json:"registered_number"`
				StateInscription string `json:"state_inscription"`
			} `json:"carrier"`
			CarrierOriginalDeliveryTime struct {
				Days          int    `json:"days"`
				EstimatedDate string `json:"estimated_date"`
			} `json:"carrier_original_delivery_time"`
			CostPrice    float64 `json:"cost_price"`
			DeliveryTime struct {
				Days          int    `json:"days"`
				EstimatedDate string `json:"estimated_date"`
			} `json:"delivery_time"`
			Esg struct {
				Co2EmissionEstimate float64 `json:"co2_emission_estimate"`
			} `json:"esg"`
			Expiration           time.Time `json:"expiration"`
			FinalPrice           float64   `json:"final_price"`
			HomeDelivery         bool      `json:"home_delivery"`
			Modal                string    `json:"modal"`
			Offer                int       `json:"offer"`
			OriginalDeliveryTime struct {
				Days          int    `json:"days"`
				EstimatedDate string `json:"estimated_date"`
			} `json:"original_delivery_time"`
			Service        string `json:"service"`
			SimulationType int    `json:"simulation_type"`
			TableReference string `json:"table_reference"`
			Weights        struct {
				Cubed float64 `json:"cubed"`
				Real  int     `json:"real"`
				Used  float64 `json:"used"`
			} `json:"weights"`
		} `json:"offers"`
		RegisteredNumberDispatcher string `json:"registered_number_dispatcher"`
		RegisteredNumberShipper    string `json:"registered_number_shipper"`
		RequestID                  string `json:"request_id"`
		ZipcodeOrigin              int    `json:"zipcode_origin"`
	} `json:"dispatchers"`
}

func NewFreteRapidoService(url string) *Service {
	var client = http.Client{}

	return &Service{
		url,
		client,
	}
}

func (s *Service) Quote(qr *QuoteSimulationRequest) (*QuoteSimulationResponse, error) {
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

	var jsonResponse = QuoteSimulationResponse{}
	json.Unmarshal(bodyResponse, &jsonResponse)

	return &jsonResponse, nil
}
