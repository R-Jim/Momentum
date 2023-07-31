package operator

import "github.com/R-jim/Momentum/aggregate/event"

func appendEvent[S event.SampleStore | event.MioStore | event.ProductStore | event.WorkerStore | event.BuildingStore | event.StageStore | event.StreetStore](store *S, e event.Event) error {
	s := event.Store(*store)
	return s.AppendEvent(e)
}
