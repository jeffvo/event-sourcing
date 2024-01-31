package dto

import (
	"bytes"
	"encoding/json"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/gofrs/uuid"
	"github.com/jeffvo/event-sourcing/internal/models"
)

type NewProductStockDTO struct {
	ID     uuid.UUID `json:"id" binding:"required"`
	Amount int64     `json:"amount" binding:"required"`
	Action models.Action
}

func ToEvent(product *NewProductStockDTO) *esdb.EventData {

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(product)

	return &esdb.EventData{
		EventID:     product.ID,
		EventType:   string(product.Action),
		ContentType: esdb.JsonContentType,
		Data:        reqBodyBytes.Bytes(),
		Metadata:    nil,
	}
}
