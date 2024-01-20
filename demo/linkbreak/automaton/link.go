package automaton

import (
	"github.com/R-jim/Momentum/demo/linkbreak/link"
	"github.com/R-jim/Momentum/demo/linkbreak/position"
	"github.com/R-jim/Momentum/demo/linkbreak/runner"
	"github.com/R-jim/Momentum/math"
	"github.com/R-jim/Momentum/template/event"
	"github.com/google/uuid"
)

type LinkAutomaton struct {
	playerID uuid.UUID

	runnerStore   *event.Store
	positionStore *event.Store
	linkStore     *event.Store

	linkOperator link.Operator

	createdLinks map[uuid.UUID][]uuid.UUID
}

func NewLinkAutomaton(playerID uuid.UUID, runnerStore, positionStore, linkStore *event.Store) LinkAutomaton {
	return LinkAutomaton{
		playerID: playerID,

		runnerStore:   runnerStore,
		positionStore: positionStore,
		linkStore:     linkStore,

		linkOperator: link.Operator{
			LinkStore: linkStore,
		},

		createdLinks: map[uuid.UUID][]uuid.UUID{},
	}
}

func (l *LinkAutomaton) CreateLinks(linkRange float64) error {
	linkableRunners, err := l.getLinkableRunner(linkRange)
	if err != nil {
		return err
	}

	result := map[uuid.UUID][]uuid.UUID{}
	for _, targetID := range linkableRunners {
		sourceID := l.playerID

		if l.isLinkExist(sourceID, targetID) {
			result[sourceID] = append(result[sourceID], targetID)
			continue
		}

		if err := l.linkOperator.NewLink(uuid.New(), sourceID, targetID); err != nil {
			return err
		}
		result[sourceID] = append(result[sourceID], targetID)
	}

	l.createdLinks = result
	return nil
}

func (l *LinkAutomaton) isLinkExist(sourceID, targetID uuid.UUID) bool {
	for existedSourceID, existedTargetIDs := range l.createdLinks {
		if existedSourceID == l.playerID {
			for _, existedTargetID := range existedTargetIDs {
				if existedTargetID == targetID {
					return true
				}
			}
		}
	}
	return false
}

func (l LinkAutomaton) getLinkableRunner(linkRange float64) ([]uuid.UUID, error) {
	result := []uuid.UUID{}

	runnerWithPositionSet, err := l.getRunnerAndPositionSet()
	if err != nil {
		return nil, err
	}
	playerPosition := runnerWithPositionSet[l.playerID]
	delete(runnerWithPositionSet, l.playerID)

	for runnerID, runnerPosition := range runnerWithPositionSet {
		_, _, distanceSqrt := math.GetDistances(playerPosition, runnerPosition)
		if distanceSqrt <= linkRange {
			result = append(result, runnerID)
		}
	}

	return result, nil
}

func (l LinkAutomaton) getRunnerAndPositionSet() (map[uuid.UUID]math.Pos, error) {
	result := map[uuid.UUID]math.Pos{}

	runnerEventMap := l.runnerStore.GetEvents()
	for runnerID, runnerEvents := range runnerEventMap {
		positionID, err := runner.GetRunnerPositionID(runnerEvents)
		if err != nil {
			return nil, err
		}

		positionEvents, err := l.positionStore.GetEventsByEntityID(positionID)
		if err != nil {
			return nil, err
		}

		position, err := position.GetPositionProjection(positionEvents)
		if err != nil {
			return nil, err
		}

		result[runnerID] = position
	}

	return result, nil
}

func (l LinkAutomaton) DeleteLinks(linkRange float64) error {
	type targetLink struct {
		id       uuid.UUID
		targetID uuid.UUID
	}

	aliveLinkSourceTargetsSet := map[uuid.UUID][]targetLink{}

	for linkID, linkEvents := range l.linkStore.GetEvents() {
		if aliveLink, err := link.GetAliveLink(linkEvents); err != nil {
			return err
		} else if aliveLink != nil {
			aliveLinkSourceTargetsSet[aliveLink.Source()] = append(aliveLinkSourceTargetsSet[aliveLink.Source()], targetLink{
				id:       linkID,
				targetID: aliveLink.Target(),
			})
		}
	}

	destroyLinkIDs := []uuid.UUID{}
	for sourceID, targetLinks := range aliveLinkSourceTargetsSet {
		sourcePosition, err := runner.GetRunnerPosition(sourceID, *l.runnerStore, *l.positionStore)
		if err != nil {
			return nil
		}

		for _, targetLink := range targetLinks {
			targetPosition, err := runner.GetRunnerPosition(targetLink.targetID, *l.runnerStore, *l.positionStore)
			if err != nil {
				return nil
			}

			if _, _, distanceSqrt := math.GetDistances(sourcePosition, targetPosition); distanceSqrt > linkRange {
				destroyLinkIDs = append(destroyLinkIDs, targetLink.id)
			}
		}
	}

	for _, linkID := range destroyLinkIDs {
		if err := l.linkOperator.Destroy(linkID); err != nil {
			return err
		}
	}
	return nil
}
