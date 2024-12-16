package conc

import (
	"github.com/BaldiSlayer/rprp/lab5/internal/defaults"
	"github.com/BaldiSlayer/rprp/lab5/internal/models"
	"sync"
	"sync/atomic"
)

var (
	barrier          = sync.NewCond(&sync.Mutex{})
	completedThreads = atomic.Int32{}
)

var requestPool = sync.Pool{
	New: func() interface{} {
		return &models.ElemRequest{
			Result: make(chan int),
		}
	},
}

type RowWorker struct {
	threadChannels []chan models.ElemRequest
	matrix         [][]int
	newMatrix      [][]int
	threadID       int
	numThreads     int
}

func (r *RowWorker) Wait() {
	barrier.L.Lock()

	completedThreads.Add(1)

	if completedThreads.Load() == int32(r.numThreads) {
		completedThreads.Store(0)

		for i := range r.matrix {
			r.matrix[i] = r.newMatrix[i]
		}

		barrier.Broadcast()
	} else {
		barrier.Wait()
	}

	barrier.L.Unlock()
}

func (r *RowWorker) Run(start, end int) {
	listenChan := r.threadChannels[r.threadID]

	go func() {
		for req := range listenChan {
			req.Result <- r.matrix[req.Y][req.X]
		}
	}()

	for iter := 0; iter < defaults.Iterations; iter++ {
		for i := start; i < end; i++ {
			for j := 0; j < defaults.Cols; j++ {
				r.newMatrix[i][j] = r.updateCellParallel(
					r.matrix,
					i, j,
					defaults.Rows,
					defaults.Cols,
					r.numThreads,
				)
			}
		}

		r.Wait()
	}

	close(r.threadChannels[r.threadID])
}

func RunParallel(matrix, newMatrix [][]int, numThreads int) {
	var wg sync.WaitGroup

	threadChannels := make([]chan models.ElemRequest, numThreads)
	for i := range threadChannels {
		threadChannels[i] = make(chan models.ElemRequest, 1)
	}

	wg.Add(numThreads)

	for t := 0; t < numThreads; t++ {
		startRow := t * defaults.Rows / numThreads
		endRow := (t + 1) * defaults.Rows / numThreads
		threadID := t

		go func() {
			defer wg.Done()

			worker := RowWorker{
				threadChannels: threadChannels,
				matrix:         matrix,
				newMatrix:      newMatrix,
				threadID:       threadID,
				numThreads:     numThreads,
			}

			worker.Run(startRow, endRow)
		}()
	}

	wg.Wait()
}

func getLiveness(threadChan chan models.ElemRequest, nx, ny int) int {
	res := make(chan int, 1)

	threadChan <- models.ElemRequest{
		X:      nx,
		Y:      ny,
		Result: res,
	}

	return <-res
}

func getLivenessWithSp(threadChan chan models.ElemRequest, nx, ny int) int {
	req := requestPool.Get().(*models.ElemRequest)
	defer requestPool.Put(req)

	req.X = nx
	req.Y = ny

	threadChan <- *req

	return <-req.Result
}

func (r *RowWorker) updateCellParallel(matrix [][]int, x, y, rowCount, colCount int, numThreads int) int {
	liveNeighbors := 0

	for _, d := range defaults.Directions {
		// уже и так умрет от перенаселения
		if liveNeighbors > 3 {
			break
		}

		nx, ny := (x+d.DX+rowCount)%rowCount, (y+d.DY+colCount)%colCount

		ownerThread := nx * numThreads / rowCount

		if ownerThread != x*numThreads/rowCount {
			liveNeighbors += getLivenessWithSp(r.threadChannels[ownerThread], nx, ny)

			continue
		}

		liveNeighbors += matrix[nx][ny]
	}

	// считаем значение
	if matrix[x][y] == 1 {
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
