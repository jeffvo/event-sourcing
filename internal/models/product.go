package models

import "github.com/gofrs/uuid"

type ProductStock struct {
	Amount int64
}

type ProductStockId struct {
	ID uuid.UUID
}
