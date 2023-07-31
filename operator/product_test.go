package operator

import (
	"testing"

	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func Test_Product_Init(t *testing.T) {
	productID := uuid.New()
	productStore := event.NewProductStore()

	productOperator := ProductOperator{
		productStore: &productStore,
	}

	err := productOperator.Init(productID, "TEST")
	require.NoError(t, err)

	events, err := event.Store(productStore).GetEventsByEntityID(productID)
	require.NoError(t, err)

	productState, err := aggregator.GetProductState(events)
	require.NoError(t, err)

	require.NotEqual(t, uuid.Nil, productState.ID)
}

func Test_Product_Progress(t *testing.T) {
	productID := uuid.New()
	productStore := event.NewProductStore()

	productOperator := ProductOperator{
		productStore: &productStore,
	}

	require.NoError(t, productOperator.Init(productID, "TEST"))

	events, err := event.Store(productStore).GetEventsByEntityID(productID)
	require.NoError(t, err)

	productState, err := aggregator.GetProductState(events)
	require.NoError(t, err)

	require.NotEqual(t, uuid.Nil, productState.ID)

	require.NoError(t, productOperator.Progress(productID, 50))

	events, err = event.Store(productStore).GetEventsByEntityID(productID)
	require.NoError(t, err)

	productState, err = aggregator.GetProductState(events)
	require.NoError(t, err)

	require.Equal(t, float64(50), productState.Progress)
}
