package app

import (
	"context"
	"github.com/BaldiSlayer/rprp/lab4/internal/fork"
	"github.com/BaldiSlayer/rprp/lab4/internal/phr"
	"sync"
)

type App struct{}

func New() *App {
	return &App{}
}

func philosophersLive(ctx context.Context, phs []phr.Philosopher) chan struct{} {
	ch := make(chan struct{})

	go func() {
		var wg sync.WaitGroup

		wg.Add(len(phs))

		for _, philosopher := range phs {
			philosopher := philosopher

			go func() {
				defer wg.Done()

				philosopher.Live(ctx)
			}()
		}

		wg.Wait()

		ch <- struct{}{}
	}()

	return ch
}

func (app *App) Run(ctx context.Context, phNum int) error {
	forks := make([]fork.Fork, phNum)

	for i := 0; i < phNum; i++ {
		forks[i] = fork.New(&sync.Mutex{})
	}

	phs := make([]phr.Philosopher, phNum)

	for i := 0; i < phNum; i++ {
		phs[i] = phr.New(
			i,
			forks[i],
			forks[(i+1)%phNum],
		)
	}

	select {
	case <-philosophersLive(ctx, phs):
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}
