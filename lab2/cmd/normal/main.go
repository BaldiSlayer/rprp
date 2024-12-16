package main

import (
	"context"
	"log"
	"time"

	"github.com/BaldiSlayer/rprp/lab2/internal/generator"
	"github.com/BaldiSlayer/rprp/lab2/internal/iterations"
	"gonum.org/v1/gonum/mat"
)

const (
	size    = 10000
	epsilon = 1e-5
)

func main() {
	matrixA := generator.GenMatrix(size)
	bData := generator.GenVectorConstN(size)

	b := mat.NewVecDense(size, bData)

	xData := generator.GenRandomVector(size)
	x := mat.NewVecDense(size, xData)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	startTime := time.Now()

	result, err := iterations.SimpleIteration(ctx, matrixA, b, x, epsilon)
	if err != nil {
		log.Fatalf("ошибка SimpleIteration: %e", err)
	}

	log.Println("Program time:", time.Since(startTime))

	printResult(matrixA, result, b, x, epsilon)
}

func printResult(A *mat.Dense, result, b, x *mat.VecDense, epsilon float64) {
	log.Println("Матрица A:\n", mat.Formatted(A))
	log.Println("Вектор b:", b)
	log.Println("Начальное приближение x:", x)
	log.Println("Желаемая точность:", epsilon)
	log.Println("Результат:\n", mat.Formatted(result))
}
