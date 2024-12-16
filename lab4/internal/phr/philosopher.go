// Package phr - philosopher package
package phr

import (
	"context"
	"github.com/BaldiSlayer/rprp/lab4/internal/fork"
	"log"
	"math/rand"
	"time"
)

func getMillis(a, b int64) time.Duration {
	randomMcs := rand.Int63n(b-a+1) + a

	return time.Duration(randomMcs) * time.Millisecond
}

type Philosopher struct {
	id        int
	leftFork  fork.Fork
	rightFork fork.Fork
}

func New(id int, leftFork, rightFork fork.Fork) Philosopher {
	return Philosopher{
		id:        id,
		leftFork:  leftFork,
		rightFork: rightFork,
	}
}

func (p *Philosopher) think() {
	log.Printf("Philosopher %d is thinking\n", p.id)
	time.Sleep(getMillis(100, 1000))
}

func (p *Philosopher) eat() {
	log.Printf("Philosopher %d is eating\n", p.id)
	time.Sleep(getMillis(100, 1000))
	log.Printf("Philosopher %d ended eating\n", p.id)
}

func (p *Philosopher) getForksIter() bool {
	if ok := p.leftFork.Get(); !ok {
		return false
	}

	if ok := p.rightFork.Get(); !ok {
		p.leftFork.Release()

		return false
	}

	return true
}

func (p *Philosopher) getForks() {
	for {
		if ok := p.getForksIter(); ok {
			return
		}

		time.Sleep(getMillis(100, 500))
	}
}

func (p *Philosopher) putAllForks() {
	p.leftFork.Release()
	p.rightFork.Release()
}

func (p *Philosopher) Live(ctx context.Context) {
	p.think()
	p.getForks()
	p.eat()
	p.putAllForks()
}
