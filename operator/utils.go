package operator

import "github.com/R-jim/Momentum/aggregate/event"

func appendEvent[S event.SampleStore](store *S, e event.Event) error {
	s := event.Store(*store)
	return s.AppendEvent(e)
}
