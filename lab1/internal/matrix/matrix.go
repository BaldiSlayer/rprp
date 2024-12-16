package matrix

import "C"
import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Matrix struct {
	data [][]int
	rows int
	cols int
}

var ErrCantMultiply = fmt.Errorf("failed to multiply matrix: the dimensions do not match")
var ErrBlocks = fmt.Errorf("failed to multiply matrix: blocks can't divide matrix normally")

// NewMatrix создает новую пустую матрицу заданных размеров
func NewMatrix(rows, cols int) Matrix {
	data := make([][]int, rows)
	for i := range data {
		data[i] = make([]int, cols)
	}
	return Matrix{data: data, rows: rows, cols: cols}
}

// NewRandomMatrix создает новую матрицу заполненную случайными числами
// вы можете добавить параметры для определения размера матрицы и диапазона случайных чисел
func NewRandomMatrix(rows, cols, maxValue int) Matrix {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	m := Matrix{
		data: make([][]int, rows),
		rows: rows,
		cols: cols,
	}

	for i := range m.data {
		m.data[i] = make([]int, cols)
		for j := range m.data[i] {
			m.data[i][j] = rng.Intn(maxValue) // случайные числа от 0 до maxValue-1
		}
	}

	return m
}

// String реализует способ печати матрицы
func (m Matrix) String() string {
	var result string
	for _, row := range m.data {
		result += fmt.Sprintln(row)
	}
	return result
}

// MultiplyMatrices перемножает две матрицы
func MultiplyMatrices(a, b *Matrix) (Matrix, error) {
	if a.cols != b.rows {
		return Matrix{}, ErrCantMultiply
	}

	result := NewMatrix(a.rows, b.cols)

	for i := 0; i < result.rows; i++ {
		for j := 0; j < result.cols; j++ {
			sum := 0
			for k := 0; k < a.cols; k++ {
				sum += a.data[i][k] * b.data[k][j]
			}
			result.data[i][j] = sum
		}
	}

	return result, nil
}

func multiplyWorker(m1, m2, result *Matrix, row, col int) {
	sum := 0
	for k := 0; k < m1.cols; k++ {
		sum += m1.data[row][k] * m2.data[k][col]
	}

	result.data[row][col] = sum
}

func MultiplyMatricesByRows(m1, m2 *Matrix, maxGoroutines int) (Matrix, error) {
	if m1.cols != m2.rows {
		return Matrix{}, ErrCantMultiply
	}

	semaphore := make(chan struct{}, maxGoroutines)

	result := NewMatrix(m1.rows, m2.cols)

	var wg sync.WaitGroup
	wg.Add(result.rows)

	for i := 0; i < result.rows; i++ {
		semaphore <- struct{}{}

		go func(i int) {
			defer wg.Done()
			defer func() { <-semaphore }()

			for j := 0; j < result.cols; j++ {
				multiplyWorker(m1, m2, &result, i, j)
			}
		}(i)
	}

	wg.Wait()

	return result, nil
}

func MultiplyMatricesByCols(m1, m2 *Matrix, maxGoroutines int) (Matrix, error) {
	if m1.cols != m2.rows {
		return Matrix{}, ErrCantMultiply
	}

	semaphore := make(chan struct{}, maxGoroutines)

	result := NewMatrix(m1.rows, m2.cols)

	var wg sync.WaitGroup
	wg.Add(result.rows)

	for i := 0; i < result.cols; i++ {
		semaphore <- struct{}{}

		go func(i int) {
			defer wg.Done()
			defer func() { <-semaphore }()

			for j := 0; j < result.rows; j++ {
				multiplyWorker(m1, m2, &result, i, j)
			}
		}(i)
	}

	wg.Wait()

	return result, nil
}

func abc(m1, m2, result *Matrix, lt, lb, rt, rb int) {
	for i := lt; i < lb; i++ {
		for j := rt; j < rb; j++ {
			//
			sum := 0
			for k := 0; k < m1.cols; k++ {
				sum += m1.data[i][k] * m2.data[k][j]
			}
			result.data[i][j] = sum
			//
		}
	}
}

// MultiplyMatricesWithBlocks перемножает матрицы с разделением на блоки.
func MultiplyMatricesWithBlocks(m1, m2 *Matrix, maxGoroutines int, w, h int) (Matrix, error) {
	if m1.cols != m2.rows {
		return Matrix{}, ErrCantMultiply
	}

	if m1.rows%w != 0 || m2.cols%h != 0 {
		return Matrix{}, ErrBlocks
	}

	result := NewMatrix(m1.rows, m2.cols)

	semaphore := make(chan struct{}, maxGoroutines)

	var wg sync.WaitGroup
	wg.Add((m1.rows / w) * (m2.cols / h))

	for i := 0; i < m1.rows/w; i++ {
		for j := 0; j < m2.cols/h; j++ {
			semaphore <- struct{}{}

			go func(i, j int) {
				defer wg.Done()
				defer func() { <-semaphore }()

				abc(m1, m2, &result, i*w, i*w+w, j*h, j*h+h)
			}(i, j)
		}
	}

	wg.Wait()

	return result, nil
}
