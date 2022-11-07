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

	if combatState.Status != jet.FlyingStatus || fuelStorage.Quantity <= 0 || combatState.Target != nil {
		return nil
	}
	fmt.Printf("[AUTO_FLY] %v\n", id)
	return i.autoPatrol(id)
}

func (i impl) autoPatrol(id string) error {
	radius := float64(40)
	targetX := float64(150)
	targetY := float64(150)

	positionState, err := jet.GetPositionState(i.jetStore, id)
	if err != nil {
		return err
	}

	x, y := getNextStepXY(positionState, targetX, targetY, radius, 2)
	inventoryState, err := jet.GetInventoryState(i.jetStore, id)
	if err != nil {
		return err
	}

	fmt.Printf("[AUTO_PATROL] %v:%v\n", x, y)
	return i.operator.Jet.Fly(id, inventoryState.FuelTankID, 1, jet.PositionState{
		X: x,
		Y: y,
	})
}
