package app

import (
	"log"
	"math/rand"
	"sync"

	"github.com/BaldiSlayer/rprp/lab52/internal/defaults"
	"github.com/BaldiSlayer/rprp/lab52/internal/llist"
)

type App struct{}

func New() *App {
	return &App{}
}

func Filling(list *llist.LinkedList) {
	wg := sync.WaitGroup{}

	wg.Add(defaults.NumThreads)

	for i := 0; i < defaults.NumThreads; i++ {
		go func() {
			defer wg.Done()

			for j := 0; j < defaults.NumValuesPerThread; j++ {
				list.AddIfNotExists(rand.Intn(defaults.ValueMaxRand))
			}
		}()
	}

	wg.Wait()
}

func (app *App) Run() {
	list := llist.New()

	Filling(list)

	log.Println(list)

	if list.CheckForDuplicates() {
		log.Println("Дубликаты найдены в списке!")

		return
	}

	log.Println("Дубликатов не найдено!!!")
}
