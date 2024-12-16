package fork

import (
	"math/rand"
	"sync"
	"time"
)

func getMillis(a, b int64) time.Duration {
	randomMcs := rand.Int63n(b-a+1) + a

	return time.Duration(randomMcs) * time.Millisecond
}

type Fork struct {
	fork *sync.Mutex
}

func New(f *sync.Mutex) Fork {
	return Fork{
		fork: f,
	}
}

func (f *Fork) Get() bool {
	if f.fork.TryLock() {
		// мы не можем взять вилку мгновенно
		time.Sleep(getMillis(100, 1000))

		return true
	}

	return false
}

func (f *Fork) Release() {
	f.fork.Unlock()
}
