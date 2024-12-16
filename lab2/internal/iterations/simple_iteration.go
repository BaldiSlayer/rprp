package iterations

import (
	"context"
	"errors"
	"log"
	"math"

	"gonum.org/v1/gonum/mat"
)

func SimpleIteration(ctx context.Context, A *mat.Dense, b, x0 *mat.VecDense, epsilon float64) (*mat.VecDense, error) {
	n, _ := A.Dims()
	x := mat.NewVecDense(n, nil)

	t := getTau(n)
	prevNormDiff := 0.0

	var iteration int
	for {
		select {
		case <-ctx.Done():
			return nil, errors.New("context deadline exceeded")
		default:
		}

		calculateMatrix(A, x, x0, b, t)
		iteration++

		if completionCriterion(x, x0, b, n, iteration, epsilon, &prevNormDiff) {
			return x, nil
		}

		x0.CopyVec(x)
	}
}

func calculateMatrix(A *mat.Dense, x, x0, b *mat.VecDense, t float64) {
	// Compare A * x0
	x.MulVec(A, x0)

	// Compare A * x0 - b
	x.SubVec(x, b)

	// Compare t * (A * x0 - b)
	x.ScaleVec(t, x)

	// Compare x0 + t * (A * x0 - b)
	x.AddVec(x0, x)
}

func completionCriterion(x, x0, b *mat.VecDense, n, i int, epsilon float64, prevNormDiff *float64) bool {
	diff := mat.NewVecDense(n, nil)
	diff.SubVec(x, x0)
	normDiff := mat.Norm(diff, 2)
	normB := mat.Norm(b, 2)
	if normDiff/normB < epsilon {
		log.Printf("Завершено после %d итераций\n", i)
		return true
	}

	if *prevNormDiff == 0 {
		*prevNormDiff = normDiff
	} else {
		// Если изменение 𝜀 меньше определенного порога (например, 1e-8), то выход
		if math.Abs(normDiff/normB-(*prevNormDiff)/normB) < 1e-8 {
			log.Printf("Завершено после %d итераций\n", i)
			return true
		}
		*prevNormDiff = normDiff
	}

	return false
}

func getTau(N int) float64 {
	return -0.1 / float64(N)
}
