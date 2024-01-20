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

func (l *LinkAutomaton) CreateOrStrengthenLinks(linkRange float64) error {
	linkableRunners, err := l.getLinkableRunner(linkRange)
	if err != nil {
		return err
	}

	activeLinkProjections := []link.LinkProjection{}
	for _, linkEvents := range l.linkStore.GetEvents() {
		linkProjection, err := link.GetLinkProjection(linkEvents)
		if err != nil {
			return err
		}

		if !linkProjection.IsDestroyed {
			activeLinkProjections = append(activeLinkProjections, linkProjection)
		}
	}

	for _, targetID := range linkableRunners {
		sourceID := l.playerID

		if err := createOrStrengthenLink(l.linkOperator, activeLinkProjections, sourceID, targetID); err != nil {
			return err
		}

	}

	return nil
}

func createOrStrengthenLink(linkOperator link.Operator, activeLinkProjections []link.LinkProjection, sourceID, targetID uuid.UUID) error {
	for _, linkProjection := range activeLinkProjections {
		if linkProjection.SourceID == sourceID && linkProjection.TargetID == targetID {
			if err := linkOperator.Strengthen(linkProjection.ID); err != nil {
				return err
			}
			return nil
		}
	}

	if err := linkOperator.New(uuid.New(), sourceID, targetID); err != nil {
		return err
	}
	return nil
}

func (l LinkAutomaton) getLinkableRunner(linkRange float64) ([]uuid.UUID, error) {
	result := []uuid.UUID{}

	runnerWithPositionSet, err := getRunnerAndPositionSet(*l.runnerStore, *l.positionStore)
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

func getRunnerAndPositionSet(runnerStore, positionStore event.Store) (map[uuid.UUID]math.Pos, error) {
	result := map[uuid.UUID]math.Pos{}

	runnerEventMap := runnerStore.GetEvents()
	for runnerID, runnerEvents := range runnerEventMap {
		if isRunnerDestroyed(runnerEvents) {
			continue
		}

		positionID, err := runner.GetRunnerPositionID(runnerEvents)
		if err != nil {
			return nil, err
		}

		positionEvents, err := positionStore.GetEventsByEntityID(positionID)
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

type DestroyLinkAutomaton struct {
	runnerStore *event.Store
	linkStore   *event.Store

	linkOperator link.Operator
}

func NewDestroyLinkAutomaton(runnerStore, linkStore *event.Store) DestroyLinkAutomaton {
	return DestroyLinkAutomaton{
		runnerStore: runnerStore,
		linkStore:   linkStore,

		linkOperator: link.Operator{
			LinkStore: linkStore,
		},
	}
}

func (d DestroyLinkAutomaton) DestroyLinkWithDestroyedRunner() error {
	activeLinks := []link.LinkProjection{}

	for _, linkEvents := range d.linkStore.GetEvents() {
		if linkProjection, err := link.GetLinkProjection(linkEvents); err != nil {
			return err
		} else if !linkProjection.IsDestroyed {
			activeLinks = append(activeLinks, linkProjection)
		}
	}

	for _, link := range activeLinks {
		sourceEvents, err := d.runnerStore.GetEventsByEntityID(link.SourceID)
		if err != nil {
			return err
		}
		targetEvents, err := d.runnerStore.GetEventsByEntityID(link.TargetID)
		if err != nil {
			return err
		}

		if isRunnerDestroyed(sourceEvents) || isRunnerDestroyed(targetEvents) {
			if err := d.linkOperator.Destroy(link.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func isRunnerDestroyed(runnerEvents []event.Event) bool {
	for i := len(runnerEvents) - 1; i >= 0; i-- {
		if runnerEvents[i].Effect == runner.DestroyEffect {
			return true
		}
	}

	return false
}
