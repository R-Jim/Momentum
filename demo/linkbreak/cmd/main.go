package main

import (
	"image/color"
	"log"

	"github.com/R-jim/Momentum/demo/linkbreak/automaton"
	"github.com/R-jim/Momentum/demo/linkbreak/runner"
	"github.com/R-jim/Momentum/math"
	"github.com/R-jim/Momentum/template/event"
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	WINDOW_X = 400
	WINDOW_Y = 300

	POINT_WIDTH  = 5
	POINT_HEIGHT = 5
)

type Game struct {
	ObjectPos math.Pos
	TargetPos math.Pos

	RunnerOperator runner.Operator

	RunnerProjector runner.Projector

	LinkAutomaton          automaton.LinkAutomaton
	BreakAutomaton         automaton.BreakAutomaton
	DestroyRunnerAutomaton automaton.DestroyRunnerAutomaton
	DestroyLinkAutomaton   automaton.DestroyLinkAutomaton

	PlayerID uuid.UUID

	Counter int
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.TargetPos = math.NewPos(float64(x), float64(y))
	}

	if !g.ObjectPos.IsEqualRound(g.TargetPos) {
		g.ObjectPos = math.NewPos(math.GetNextStepXY(g.ObjectPos, 0, g.TargetPos, 0, 1, 180))

		runner, err := g.RunnerProjector.GetRunnerProjection(g.PlayerID)
		if err != nil {
			return err
		}
		g.RunnerOperator.MoveRunner(runner.ID, g.ObjectPos)
	}

	if g.Counter%60 == 0 {
		if err := g.LinkAutomaton.CreateOrStrengthenLinks(50); err != nil {
			return err
		}
		if err := g.LinkAutomaton.DeleteLinks(50); err != nil {
			return err
		}
		if err := g.BreakAutomaton.BreakLinkedRunners(3); err != nil {
			return err
		}
		if err := g.DestroyRunnerAutomaton.DestroyEmptyHealthRunner(); err != nil {
			return err
		}
		if err := g.DestroyLinkAutomaton.DestroyLinkWithDestroyedRunner(); err != nil {
			return err
		}
	}
	}

	g.Counter++

	return nil
}

func (g *Game) SpawnEntity(id uuid.UUID, factionValue int, position math.Pos) error {
	if err := g.RunnerOperator.NewRunner(id, 5, factionValue, position); err != nil {
		return err
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	runnerProjections, err := g.RunnerProjector.GetRunnerProjections()
	if err != nil {
		log.Fatal(err)
	}

	for _, runnerProjection := range runnerProjections {
		if runnerProjection.IsDestroyed {
			continue
		}

		var clr color.Color
		switch runnerProjection.Faction {
		case 1:
			clr = color.RGBA{0x0, 0xff, 0x0, 0xff}
		case 2:
			clr = color.RGBA{0xff, 0x0, 0x0, 0xff}
		default:
			clr = color.White
		}

		ebitenutil.DrawRect(screen, runnerProjection.Position.X-POINT_WIDTH/2, runnerProjection.Position.Y-POINT_HEIGHT/2, POINT_WIDTH, POINT_HEIGHT, clr)
	}

	linkProjections, err := g.RunnerProjector.GetLinkProjections()
	if err != nil {
		log.Fatal(err)
	}
	for _, linkProjection := range linkProjections {
		clr := color.RGBA{0x0, 0xff, 0x0, 0xff}
		ebitenutil.DrawLine(screen, linkProjection.OwnerPosition.X, linkProjection.OwnerPosition.Y, linkProjection.TargetPosition.X, linkProjection.TargetPosition.Y, clr)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(WINDOW_X, WINDOW_Y)
	ebiten.SetWindowTitle("Link break")

	g := &Game{}
	g.Init()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Init() {
	runnerStore := event.NewStore()
	healthStore := event.NewStore()
	positionStore := event.NewStore()
	linkStore := event.NewStore()

	g.RunnerOperator = runner.Operator{
		RunnerStore:   &runnerStore,
		HealthStore:   &healthStore,
		PositionStore: &positionStore,
	}

	g.RunnerProjector = runner.Projector{
		RunnerStore:   &runnerStore,
		HealthStore:   &healthStore,
		PositionStore: &positionStore,
		LinkStore:     &linkStore,
	}

	g.ObjectPos = math.NewPos(WINDOW_X/2, WINDOW_Y/2)
	playerID := uuid.New()
	if err := g.SpawnEntity(playerID, 1, g.ObjectPos); err != nil {
		log.Fatal(err)
	}
	g.PlayerID = playerID

	g.LinkAutomaton = automaton.NewLinkAutomaton(playerID, &runnerStore, &positionStore, &linkStore)
	g.BreakAutomaton = automaton.NewBreakAutomaton(&linkStore, &runnerStore, &healthStore)
	g.DestroyRunnerAutomaton = automaton.NewDestroyRunnerAutomaton(&runnerStore, &positionStore, &healthStore)
	g.DestroyLinkAutomaton = automaton.NewDestroyLinkAutomaton(&runnerStore, &linkStore)

	_, err = g.SpawnEntity(2, math.NewPos(WINDOW_X/4, WINDOW_Y/4))
	if err != nil {
		log.Fatal(err)
	}

	g.TargetPos = g.ObjectPos
}
