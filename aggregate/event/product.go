package event

import "github.com/google/uuid"

const (
	ProductInitEffect Effect = "PRODUCT_INIT"

	ProductProgressEffect Effect = "PRODUCT_PROGRESS"
)

type ProductStore Store

func NewProductStore() ProductStore {
	return ProductStore(newStore())
}

func (s ProductStore) NewProductInitEvent(entityID uuid.UUID, productType string) Event {
	return Store(s).newEvent(entityID, ProductInitEffect, productType)
}

func (s ProductStore) NewProductProgressEvent(entityID uuid.UUID, value float64) Event {
	return Store(s).newEvent(entityID, ProductProgressEffect, value)
}
