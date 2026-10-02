package worker

import (
	"context"
	"sync"
)

// Group runs background goroutines (rates fetcher, etc.) with channel-based
// shutdown. Cancel the parent context or call Stop to signal all workers.
type Group struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	errCh  chan error
}

// NewGroup derives a cancellable group from parent.
func NewGroup(parent context.Context) *Group {
	ctx, cancel := context.WithCancel(parent)
	return &Group{ctx: ctx, cancel: cancel, errCh: make(chan error, 8)}
}

// Ctx is the group context workers select on.
func (g *Group) Ctx() context.Context { return g.ctx }

// Go starts fn in a goroutine; the first error is kept on Err().
func (g *Group) Go(fn func(ctx context.Context) error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		if err := fn(g.ctx); err != nil {
			select {
			case g.errCh <- err:
			default:
			}
		}
	}()
}

// Stop cancels the context and waits for workers.
func (g *Group) Stop() { g.cancel(); g.wg.Wait() }

// Err returns the first worker error, or nil.
func (g *Group) Err() error {
	select {
	case err := <-g.errCh:
		return err
	default:
		return nil
	}
}
