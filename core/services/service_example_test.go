package services_test

import (
	"context"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
)

type example struct {
	services.ServiceCtx
	g      *services.Group
	workCh chan func()
}

func (e *example) start(context.Context) error {
	e.g.Go(func() {
		for {
			select {
			case <-e.g.Closed():
			case <-time.After(time.Minute):
				e.do(e.workCh)
			}
		}
	})
	return nil
}

func (e *example) do(workCh <-chan func()) {
	// do work until none is left
	for {
		select {
		case <-e.g.Closed():
			return
		case work, ok := <-workCh:
			if !ok {
				return
			}
			work()
		default:
			return
		}
	}
}

func NewExample(lggr logger.Logger) services.ServiceCtx {
	e := &example{
		workCh: make(chan func()),
	}
	e.ServiceCtx, e.g = services.New(services.Spec{
		Name:        "Example",
		Start:       e.start,
		SubServices: nil, // optional
	}, lggr)
	return e
}

func Example() {
	NewExample(logger.NullLogger)
}
