package transformers

import (
	"strconv"

	"github.com/jsGolden/frete-rapido-api/models"
	"github.com/jsGolden/frete-rapido-api/services"
)

func TransformQuoteToFreteRapido(q *models.CreateQuoteRequest) (*services.FreteRapidoQuoteSimulationRequest, error) {
	const (
		cnpj              = "25438296000158"
		token             = "1d52a9b6b78cf07b08586152459a5c90"
		platformCode      = "5AKVkHqCn"
		dispatcherZipcode = 29161376
	)

	output := services.FreteRapidoQuoteSimulationRequest{}
	recipientZipcode, err := strconv.Atoi(q.Recipient.Address.Zipcode)
	if err != nil {
		return nil, err
	}

	output.Recipient.Zipcode = recipientZipcode
	output.Recipient.Country = "BRA"
	output.Recipient.Type = 1

	output.Shipper.RegisteredNumber = cnpj
	output.Shipper.PlatformCode = platformCode
	output.Shipper.Token = token

	dispatcher := struct {
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
	}{
		RegisteredNumber: cnpj,
		Zipcode:          dispatcherZipcode,
	}

	for _, volume := range q.Volumes {
		volumeData := struct {
			Category      string  `json:"category"`
			Amount        int     `json:"amount"`
			UnitaryWeight float64 `json:"unitary_weight"`
			UnitaryPrice  float64 `json:"unitary_price"`
			Sku           string  `json:"sku"`
			Height        float64 `json:"height"`
			Width         float64 `json:"width"`
			Length        float64 `json:"length"`
		}{
			Category:      strconv.Itoa(volume.Category),
			Amount:        volume.Amount,
			UnitaryWeight: volume.UnitWeight,
			UnitaryPrice:  volume.Price,
			Sku:           volume.SKU,
			Height:        volume.Height,
			Width:         volume.Width,
			Length:        volume.Length,
		}
		dispatcher.Volumes = append(dispatcher.Volumes, volumeData)
	}

	output.Dispatchers = append(output.Dispatchers, dispatcher)
	output.SimulationType = []int{0}
	return &output, nil
}
