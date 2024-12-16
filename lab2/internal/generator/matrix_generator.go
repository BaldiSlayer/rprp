package generator

import (
	"math/rand"

	"gonum.org/v1/gonum/mat"
)

func GenMatrix(size int) *mat.Dense {
	A := mat.NewDense(size, size, nil)

	for i := 0; i < size; i++ {
		A.Set(i, i, 2.0)

		for j := 0; j < size; j++ {
			if i != j {
				A.Set(i, j, 1.0)
			}
		}
	}

	return A
}

// GenVectorConstN создает вектор, где все значения равны n+1
func GenVectorConstN(n int) []float64 {
	bData := make([]float64, n)

	for i := 0; i < n; i++ {
		bData[i] = float64(n + 1)
	}

	return bData
}

func GenChunkMatrix(start, end, size int) *mat.Dense {
	rows := end - start
	res := mat.NewDense(rows, size, make([]float64, rows*size))

	for i := start; i < end; i++ {
		for j := 0; j < size; j++ {
			if i == j {
				res.Set(i-start, j, 2.0)
			} else {
				res.Set(i-start, j, 1.0)
			}
		}
	}

	return res
}

func GenCheckFreeVector(matrixChunkRowsSize, size int) *mat.VecDense {
	resBuf := make([]float64, matrixChunkRowsSize)

	for i := 0; i < matrixChunkRowsSize; i++ {
		resBuf[i] = float64(size + 1)
	}

	return mat.NewVecDense(matrixChunkRowsSize, resBuf)
}

func GenRandomVector(size int) []float64 {
	res := make([]float64, size)

	for i := 0; i < size; i++ {
		res[i] = float64(rand.Intn(size))
	}

	return res
}
