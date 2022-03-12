package waiter

import (
	"context"

	"golang.org/x/sync/semaphore"
)

type SigWaiter struct {
	w *semaphore.Weighted
}

func NewSigWaiter(n int64) *SigWaiter {
	return &SigWaiter{w: semaphore.NewWeighted(n)}
}

func (s *SigWaiter) Stw() {
	panic("implement me")
}

func (s *SigWaiter) Restart() {
	panic("implement me")
}

func (s *SigWaiter) Wait(ctx context.Context) {
	s.w.Acquire(ctx, 1)
}

func (s *SigWaiter) Release() {
	s.w.Release(1)
}
