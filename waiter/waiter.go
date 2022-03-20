package waiter

import (
	"context"
	"sync"

	"golang.org/x/time/rate"
)

type Interface interface {
	Wait(context.Context)
	Release()

	Stw()
	Restart()
}

type Waiter struct {
	count int
	mu    *sync.Mutex
	cond  *sync.Cond

	rt *rate.Limiter
}

func NewWaiter(r rate.Limit, b int) *Waiter {
	mu := &sync.Mutex{}

	return &Waiter{
		mu:   mu,
		cond: sync.NewCond(mu),
		rt:   rate.NewLimiter(r, b),
	}
}

func (w *Waiter) Release() {}

func (w *Waiter) Stw() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.count++
}

func (w *Waiter) Restart() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.count--
	w.cond.Broadcast()
}

func (w *Waiter) Wait(ctx context.Context) {
	w.rt.Wait(ctx)
	var rewait bool

	func() {
		w.mu.Lock()
		defer w.mu.Unlock()
		for w.count != 0 {
			w.cond.Wait()
			rewait = true
		}
	}()
	if rewait {
		w.rt.Wait(ctx)
	}
}
