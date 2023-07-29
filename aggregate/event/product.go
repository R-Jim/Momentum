package event

import "github.com/google/uuid"

const (
	ProductInitEffect Effect = "PRODUCT_INIT"

	ProductProgressEffect Effect = "PRODUCT_PROGRESS"
)

func NewProductInitEvent(entityID uuid.UUID, productType string) Event {
	return newEvent(entityID, 1, ProductInitEffect, productType)
}

func NewProductProgressEvent(entityID uuid.UUID, version int, value float64) Event {
	return newEvent(entityID, version, ProductProgressEffect, value)
}
