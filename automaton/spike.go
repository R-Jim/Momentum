package automaton

import (
	"github.com/R-jim/Momentum/aggregate/carrier"
	"github.com/R-jim/Momentum/aggregate/knight"
	"github.com/R-jim/Momentum/aggregate/spike"
	"github.com/R-jim/Momentum/operator"
	"github.com/R-jim/Momentum/util"
)

/*
TODO:
attack carrier/knight when in range 1 from spike
*/

type SpikeAutomaton interface {
	Auto(id string) error
}

type spikeImpl struct {
	carrierStore carrier.Store
	knightStore  knight.Store
	spikeStore   spike.Store

	operator operator.Operator
}

func NewSpikeAutomaton(carrierStore carrier.Store, knightStore knight.Store, spikeStore spike.Store, operator operator.Operator) SpikeAutomaton {
	return spikeImpl{
		carrierStore: carrierStore,
		knightStore:  knightStore,
		spikeStore:   spikeStore,

		operator: operator,
	}
}

func (i spikeImpl) Auto(id string) error {
	return i.AutoStrike(id)
}

func (i spikeImpl) AutoStrike(id string) error {
	positionState, err := spike.GetPositionState(i.spikeStore, id)
	if err != nil {
		return err
	}

	targetKnightIDs := []string{}
	// targetCarrierIDs := []string{}

	for _, knightID := range i.knightStore.GetEntityIDs() {
		knightCombatState, err := knight.GetState(i.knightStore, knightID)
		if err != nil {
			return err
		}
		if knightCombatState.Health.Value <= 0 {
			continue
		}

		knightPositionState, err := knight.GetPositionState(i.knightStore, knightID)
		if err != nil {
			return err
		}
		_, _, distance := util.GetDistances(positionState.X, positionState.Y, knightPositionState.X, knightPositionState.Y)
		if distance <= 1 {
			targetKnightIDs = append(targetKnightIDs, knightID)
		}
	}

	// for _, carrierID := range i.carrierStore.GetEntityIDs() {
	// 	carrierCombatState, err := carrier.GetCombatState(i.carrierStore, carrierID)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if carrierCombatState.Health.Value <= 0 {
	// 		continue
	// 	}

	// 	carrierPositionState, err := carrier.GetPositionState(i.carrierStore, carrierID)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	_, _, distance := util.GetDistances(positionState.X, positionState.Y, carrierPositionState.X, carrierPositionState.Y)
	// 	if distance <= 1 {
	// 		targetCarrierIDs = append(targetCarrierIDs, carrierID)
	// 	}
	// }

	for _, targetKnightID := range targetKnightIDs {
		err = i.operator.Spike.StrikeKnight(id, targetKnightID)
		if err != nil {
			return err
		}
	}
	// for _, targetCarrierID := range targetCarrierIDs {
	// 	err = i.operator.Spike.StrikeKnight(id, targetCarrierID)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	return nil
}
