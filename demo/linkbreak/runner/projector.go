package runner

import (
	"github.com/google/uuid"

	"github.com/R-jim/Momentum/demo/linkbreak/health"
	"github.com/R-jim/Momentum/demo/linkbreak/link"
	"github.com/R-jim/Momentum/demo/linkbreak/position"
	"github.com/R-jim/Momentum/math"
	"github.com/R-jim/Momentum/template/event"
)

type RunnerLinkProjection struct {
	Owner  uuid.UUID
	Target uuid.UUID

	OwnerPosition  math.Pos
	TargetPosition math.Pos
}

type RunnerProjection struct {
	ID                 uuid.UUID
	Faction            int
	BaseHealthValue    int
	CurrentHealthValue int
	Position           math.Pos
	IsDestroyed        bool
}

type Projector struct {
	RunnerStore   *event.Store
	HealthStore   *event.Store
	PositionStore *event.Store
	LinkStore     *event.Store
}

func (p Projector) GetRunnerProjections() ([]RunnerProjection, error) {
	projections := []RunnerProjection{}

	for runnerID := range p.RunnerStore.GetEvents() {
		projection, err := p.GetRunnerProjection(runnerID)
		if err != nil {
			return nil, err
		}

		projections = append(projections, projection)
	}

	return projections, nil
}

func (p Projector) GetRunnerProjection(id uuid.UUID) (RunnerProjection, error) {
	projection := RunnerProjection{
		ID: id,
	}

	var healthID, positionID uuid.UUID
	if events, err := p.RunnerStore.GetEventsByEntityID(id); err == nil {
		for _, e := range events {
			switch e.Effect {
			case InitEffect:
				runner, err := event.ParseData[Runner](e)
				if err != nil {
					return RunnerProjection{}, err
				}
				healthID = runner.healthID
				positionID = runner.positionID
				projection.Faction = runner.faction

			case DestroyEffect:
				projection.IsDestroyed = true
			}
		}
	} else {
		return RunnerProjection{}, err
	}

	if events, err := p.HealthStore.GetEventsByEntityID(healthID); err == nil {
		for _, e := range events {
			switch e.Effect {
			case health.InitEffect:
				health, err := event.ParseData[health.Health](e)
				if err != nil {
					return RunnerProjection{}, err
				}
				projection.BaseHealthValue = health.BaseValue()
				projection.CurrentHealthValue = health.BaseValue()
			}
		}
	} else {
		return RunnerProjection{}, err
	}

	if events, err := p.PositionStore.GetEventsByEntityID(positionID); err == nil {
		for _, e := range events {
			switch e.Effect {
			case position.InitEffect:
				positionValue, err := event.ParseData[math.Pos](e)
				if err != nil {
					return RunnerProjection{}, err
				}
				projection.Position = positionValue
			case position.MoveEffect:
				positionValue, err := event.ParseData[math.Pos](e)
				if err != nil {
					return RunnerProjection{}, err
				}
				projection.Position = positionValue
			}
		}
	} else {
		return RunnerProjection{}, err
	}

	return projection, nil
}

func (p Projector) GetLinkProjections() ([]RunnerLinkProjection, error) {
	projections := []RunnerLinkProjection{}

	aliveLinkSet, err := getLinkSet(*p.LinkStore)
	if err != nil {
		return nil, err
	}

	for sourceID, targetIDs := range aliveLinkSet {
		sourcePosition, err := GetRunnerPosition(sourceID, *p.RunnerStore, *p.PositionStore)
		if err != nil {
			return nil, err
		}

		for _, targetID := range targetIDs {
			targetPosition, err := GetRunnerPosition(targetID, *p.RunnerStore, *p.PositionStore)
			if err != nil {
				return nil, err
			}

			projections = append(projections, RunnerLinkProjection{
				Owner:  sourceID,
				Target: targetID,

				OwnerPosition:  sourcePosition,
				TargetPosition: targetPosition,
			})
		}
	}

	return projections, nil
}

func getLinkSet(linkStore event.Store) (map[uuid.UUID][]uuid.UUID, error) {
	result := map[uuid.UUID][]uuid.UUID{}

	for _, linkEvents := range linkStore.GetEvents() {
		if aliveLink, err := link.GetAliveLink(linkEvents); err != nil {
			return nil, err
		} else if aliveLink != nil {
			result[aliveLink.Source()] = append(result[aliveLink.Source()], aliveLink.Target())
		}
	}

	return result, nil
}

func GetRunnerPosition(runnerID uuid.UUID, runnerStore, positionStore event.Store) (math.Pos, error) {
	sourceEvents, err := runnerStore.GetEventsByEntityID(runnerID)
	if err != nil {
		return math.Pos{}, err
	}

	positionID, err := GetRunnerPositionID(sourceEvents)
	if err != nil {
		return math.Pos{}, err
	}

	positionEvents, err := positionStore.GetEventsByEntityID(positionID)
	if err != nil {
		return math.Pos{}, err
	}

	position, err := position.GetPositionProjection(positionEvents)
	if err != nil {
		return math.Pos{}, err
	}

	return position, nil
}

func GetRunnerPositionID(events []event.Event) (uuid.UUID, error) {
	for _, e := range events {
		if e.Effect == InitEffect {
			runner, err := event.ParseData[Runner](e)
			if err != nil {
				return uuid.UUID{}, err
			}
			return runner.PositionID(), nil
		}
	}
	return uuid.UUID{}, ErrPositionIDNotFound
}
