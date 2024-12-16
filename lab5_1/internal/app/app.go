package app

import (
	"fmt"
	"github.com/BaldiSlayer/rprp/lab5/internal/conc"
	"github.com/BaldiSlayer/rprp/lab5/internal/defaults"
	"github.com/BaldiSlayer/rprp/lab5/internal/seq"
	"github.com/BaldiSlayer/rprp/lab5/internal/tmeas"
	"math/rand"
)

type App struct{}

func New() *App {
	return &App{}
}

func (app *App) Run() {
	matrix, newMatrix := getMatrices()

	measureSequential(matrix, newMatrix)

	measureParallel(matrix, newMatrix)
}

func getMatrices() ([][]int, [][]int) {
	matrix := make([][]int, defaults.Rows)
	newMatrix := make([][]int, defaults.Rows)
	for i := range matrix {
		matrix[i] = make([]int, defaults.Cols)
		newMatrix[i] = make([]int, defaults.Cols)
		for j := range matrix[i] {
			matrix[i][j] = rand.Intn(2)
		}
	}

	return matrix, newMatrix
}

func measureSequential(matrix, newMatrix [][]int) {
	tmeas.MeasureFuncTime("Среднее время одной итерации без concurency", func() {
		seq.RunSequential(matrix, newMatrix)
	})
}

func measureParallelStep(matrix, newMatrix [][]int, numThreads int) {
	tmeas.MeasureFuncTime(
		fmt.Sprintf("Среднее время одной итерации с %d горутинами", numThreads),
		func() {
			conc.RunParallel(matrix, newMatrix, numThreads)
		},
	)
}

func measureParallel(matrix, newMatrix [][]int) {
	for _, numThreads := range defaults.Threads {
		measureParallelStep(matrix, newMatrix, numThreads)
	}
}
