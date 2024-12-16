package seq

import "github.com/BaldiSlayer/rprp/lab5/internal/defaults"

func getNewCellValue(cellValue int, liveNeighbors int) int {
	if cellValue == 1 {
		if liveNeighbors < 2 || liveNeighbors > 3 {
			return 0
		}

		return 1
	}

	if liveNeighbors == 3 {
		return 1
	}

	return 0
}

func updateCell(matrix [][]int, x, y int) int {
	liveNeighbors := 0

	for _, d := range defaults.Directions {
		nx, ny := (x+d.DX+defaults.Rows)%defaults.Rows, (y+d.DY+defaults.Cols)%defaults.Cols
		liveNeighbors += matrix[nx][ny]
	}

	return getNewCellValue(matrix[x][y], liveNeighbors)
}

func updateMatrix(matrix, newMatrix [][]int) [][]int {
	for i := 0; i < defaults.Rows; i++ {
		for j := 0; j < defaults.Cols; j++ {
			newMatrix[i][j] = updateCell(matrix, i, j)
		}
	}

	return newMatrix
}

func RunSequential(matrix, newMatrix [][]int) {
	for iter := 0; iter < defaults.Iterations; iter++ {
		newMatrix = updateMatrix(matrix, newMatrix)

		for i := range matrix {
			copy(matrix[i], newMatrix[i])
		}
	}
}
