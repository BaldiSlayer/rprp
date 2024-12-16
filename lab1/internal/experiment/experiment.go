package experiment

import (
	"errors"
	"github.com/BaldiSlayer/rprp/lab1/internal/matrix"
	"log"
	"time"
)

type Experiment struct {
	rowsA     int
	colsA     int
	maxValueA int

	rowsB     int
	colsB     int
	maxValueB int

	a matrix.Matrix
	b matrix.Matrix
}

func New(rowsA, colsA, maxValueA, rowsB, colsB, maxValueB int) *Experiment {
	return &Experiment{
		rowsA:     rowsA,
		colsA:     colsA,
		maxValueA: maxValueA,
		rowsB:     rowsB,
		colsB:     colsB,
		maxValueB: maxValueB,
		a:         matrix.NewRandomMatrix(rowsA, colsA, maxValueA),
		b:         matrix.NewRandomMatrix(rowsB, colsB, maxValueB),
	}
}

type result struct {
	m       matrix.Matrix
	elapsed time.Duration
}

// measureTime замеряет время выполнения функции
func measureTime(
	f func(a, b *matrix.Matrix, maxGorutines, w, h int) (matrix.Matrix, error),
	a, b *matrix.Matrix,
	maxGorutines, w, h int,
) (result, error) {
	start := time.Now()
	res, err := f(a, b, maxGorutines, w, h)
	elapsed := time.Since(start)

	return result{
		m:       res,
		elapsed: elapsed,
	}, err
}

type Measures struct {
	Normal time.Duration
	Cols   time.Duration
	Rows   time.Duration
	Blocks time.Duration
}

func (e *Experiment) Run(maxGorutines int) (Measures, error) {
	normal, err := measureTime(WrapMultiplyMatrices, &e.a, &e.b, 0, 0, 0)
	if err != nil {
		return Measures{}, err
	}

	log.Print("Закончен подсчет обычным способом")

	cols, err := measureTime(WrapMultiplyMatricesByCols, &e.a, &e.b, maxGorutines, 0, 0)
	if err != nil {
		return Measures{}, err
	}

	log.Print("Закончен подсчет по строкам")

	rows, err := measureTime(WrapMultiplyMatricesByRows, &e.a, &e.b, maxGorutines, 0, 0)
	if err != nil {
		return Measures{}, err
	}

	log.Print("Закончен подсчет по столбцам")

	// сделать автовыбор w и h
	blocks, err := measureTime(WrapMultiplyMatricesWithBlocks, &e.a, &e.b, maxGorutines, 10, 10)
	if err != nil {
		return Measures{}, err
	}

	// вынести в отдельную функцию
	equal := normal.m.String() == cols.m.String() && cols.m.String() == rows.m.String() && rows.m.String() == blocks.m.String()
	if !equal {
		return Measures{}, errors.New("matrices are not equal")
	}

	return Measures{
		Normal: normal.elapsed,
		Cols:   cols.elapsed,
		Rows:   rows.elapsed,
		Blocks: blocks.elapsed,
	}, nil
}

func WrapMultiplyMatrices(a, b *matrix.Matrix, _, _, _ int) (matrix.Matrix, error) {
	defer log.Print("Закончен обычный подсчет")

	log.Print("Начат обычный подсчет")

	return matrix.MultiplyMatrices(a, b)
}

func WrapMultiplyMatricesByCols(a, b *matrix.Matrix, maxGorutines, _, _ int) (matrix.Matrix, error) {
	defer log.Print("Закончен подсчет по строкам")

	log.Print("Начат подсчет по строкам")

	return matrix.MultiplyMatricesByCols(a, b, maxGorutines)
}

func WrapMultiplyMatricesByRows(a, b *matrix.Matrix, maxGorutines, _, _ int) (matrix.Matrix, error) {
	defer log.Print("Закончен подсчет по столбцам")

	log.Print("Начат подсчет по столбцам")

	return matrix.MultiplyMatricesByRows(a, b, maxGorutines)
}

func WrapMultiplyMatricesWithBlocks(a, b *matrix.Matrix, maxGorutines, w, h int) (matrix.Matrix, error) {
	defer log.Print("Закончен подсчет по блокам")

	log.Print("Начат подсчет по блокам")

	return matrix.MultiplyMatricesWithBlocks(a, b, maxGorutines, w, h)
}
