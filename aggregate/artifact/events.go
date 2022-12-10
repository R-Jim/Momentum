package artifact

type Data interface{}

type Effect string

var (
	InitEffect Effect = "ARTIFACT_INITIALED"

	SpawnSpikeEffect Effect = "ARTIFACT_SPIKE_SPAWNED"

	MoveEffect Effect = "ARTIFACT_MOVED"

	GatherEffect Effect = "ARTIFACT_GATHERED"
	DropEFfect   Effect = "ARTIFACT_DROPPED"
)

func (e Effect) IsValid() bool {
	return e == InitEffect ||
		e == SpawnSpikeEffect ||
		e == MoveEffect ||
		e == GatherEffect ||
		e == DropEFfect
}

type Event struct {
	ID     string
	Effect Effect
	Data   Data
}

func NewInitEvent(id string) Event {
	return Event{
		ID:     id,
		Effect: InitEffect,
	}
}

type spikeSpawnData struct {
	spikeID       string
	positionState PositionState
}

func NewSpawnSpikeEffect(id string, spikeID string, x, y float64) Event {
	return Event{
		ID:     id,
		Effect: SpawnSpikeEffect,
		Data: spikeSpawnData{
			spikeID: spikeID,
			positionState: PositionState{
				X: x,
				Y: y,
			},
		},
	}
}

func NewMoveEffect(id string, x, y float64) Event {
	return Event{
		ID:     id,
		Effect: MoveEffect,
		Data: PositionState{
			X: x,
			Y: y,
		},
	}
}

func NewGatherEffect(id string, gatherID string) Event {
	return Event{
		ID:     id,
		Effect: MoveEffect,
		Data:   gatherID,
	}
}

func NewDropEffect(id string, x, y float64) Event {
	return Event{
		ID:     id,
		Effect: MoveEffect,
		Data: PositionState{
			X: x,
			Y: y,
		},
	}
}
