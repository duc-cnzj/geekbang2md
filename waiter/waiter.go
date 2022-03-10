package waiter

import (
	"context"
	"log"
	"sync"

	"golang.org/x/time/rate"
)

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

func (w *Waiter) Stw() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.count++
	//w.cond.Broadcast()
	log.Println("[Stop the world]!")
}

func (w *Waiter) Restart() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.count--
	w.cond.Broadcast()
	log.Println("[Restart the world]!: ", w.count)
}

func (w *Waiter) Wait(ctx context.Context) {
	w.rt.Wait(ctx)
	var rewait bool

	func() {
		w.mu.Lock()
		defer w.mu.Unlock()
		for w.count != 0 {
			log.Println("wait")
			w.cond.Wait()
			rewait = true
		}
	}()
	if rewait {
		w.rt.Wait(ctx)
	}
}
