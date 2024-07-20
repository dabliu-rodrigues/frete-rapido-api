package transformers

import (
	"strconv"

	"github.com/jsGolden/frete-rapido-api/models"
	freterapido "github.com/jsGolden/frete-rapido-api/services/frete-rapido"
)

func TransformQuoteToFreteRapido(q *models.CreateQuoteRequest) (*freterapido.QuoteSimulationRequest, error) {
	const (
		cnpj              = "25438296000158"
		token             = "1d52a9b6b78cf07b08586152459a5c90"
		platformCode      = "5AKVkHqCn"
		dispatcherZipcode = 29161376
	)

	recipientZipcode, err := strconv.Atoi(q.Recipient.Address.Zipcode)
	if err != nil {
		return nil, err
	}

	output := freterapido.QuoteSimulationRequest{
		Recipient: freterapido.Recipient{
			Zipcode: recipientZipcode,
			Country: "BRA",
			Type:    1,
		},
		Shipper: freterapido.Shipper{
			RegisteredNumber: cnpj,
			PlatformCode:     platformCode,
			Token:            token,
		},
	}

	dispatcher := freterapido.Dispatcher{
		RegisteredNumber: cnpj,
		Zipcode:          dispatcherZipcode,
		Volumes:          make([]freterapido.Volume, len(q.Volumes)),
	}

	for i, volume := range q.Volumes {
		volumeData := freterapido.Volume{
			Category:      strconv.Itoa(volume.Category),
			Amount:        volume.Amount,
			UnitaryWeight: volume.UnitWeight,
			UnitaryPrice:  volume.Price,
			Sku:           volume.SKU,
			Height:        volume.Height,
			Width:         volume.Width,
			Length:        volume.Length,
		}
		dispatcher.Volumes[i] = volumeData
	}

	output.Dispatchers = append(output.Dispatchers, dispatcher)
	output.SimulationType = []int{0}
	return &output, nil
}
