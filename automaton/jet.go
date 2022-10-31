package automaton

import (
	"fmt"

	"github.com/R-jim/Momentum/aggregate/jet"
	"github.com/R-jim/Momentum/aggregate/storage"
	"github.com/R-jim/Momentum/operator"
)

type JetAutomaton interface {
	Auto(id string) error
}

type impl struct {
	jetStore     jet.Store
	storageStore storage.Store

	operator operator.Operator
}

func NewJetAutomaton(jetStore jet.Store, storageStore storage.Store, operator operator.Operator) JetAutomaton {
	return impl{
		jetStore:     jetStore,
		storageStore: storageStore,

		operator: operator,
	}
}

func (i impl) Auto(id string) error {
	return i.autoFly(id)
}

func (i impl) autoFly(id string) error {
	combatState, err := jet.GetCombatState(i.jetStore, id)
	if err != nil {
		return err
	}

	inventoryState, err := jet.GetInventoryState(i.jetStore, id)
	if err != nil {
		return err
	}

	fuelStorage, err := storage.GetState(i.storageStore, inventoryState.FuelTankID)
	if err != nil {
		return err
	}

	if combatState.Status != jet.FlyingStatus || fuelStorage.Quantity <= 0 || combatState.TargetID != "" {
		return nil
	}
	fmt.Printf("[AUTO_FLY] %v\n", id)

	// TODO: get calculated position
	positionState, err := jet.GetPositionState(i.jetStore, id)
	if err != nil {
		return err
	}

	return i.operator.Jet.Fly(id, inventoryState.FuelTankID, 1, jet.PositionState{
		X: positionState.X + 1,
		Y: positionState.Y + 1,
	})
}
