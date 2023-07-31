package automaton

import (
	"fmt"

	"github.com/R-jim/Momentum/aggregate/aggregator"
	"github.com/R-jim/Momentum/aggregate/event"
	"github.com/google/uuid"
)

type StageAutomaton struct {
	EntityID uuid.UUID

	ProductID uuid.UUID

	ProductStore *event.ProductStore
}

func (s StageAutomaton) IsProductFinished() bool {
	events := event.Store(*s.ProductStore).GetEvents()[s.ProductID]

	productState, err := aggregator.GetProductState(events)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return productState.IsFinish()
}
