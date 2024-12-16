package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/BaldiSlayer/rprp/lab2/internal/generator"
	"github.com/BaldiSlayer/rprp/lab2/internal/pariter"
	"github.com/emer/empi/mpi"
)

const (
	eps = 1e-5
)

var (
	size = 20000
)

func init() {
	flag.IntVar(&size, "n", size, "size of matrix")
}

func main() {
	startTime := time.Now()

	flag.Parse()

	mpi.Init()
	defer mpi.Finalize()

	comm, errC := mpi.NewComm(nil)
	if errC != nil {
		log.Fatalf("can't init mpi comm %e", errC)
	}

	start, end, buff := initMpiMatrix(comm)

	err := comm.BcastF64(0, buff)
	if err != nil {
		log.Fatalf("can't broadcast buf %e", err)
	}

	baseMatrixChunk := generator.GenChunkMatrix(start, end, size)
	freeConst := generator.GenCheckFreeVector(baseMatrixChunk.RawMatrix().Rows, size)
	solver := pariter.NewSolverWithVecSeparation(comm, baseMatrixChunk, freeConst, size, eps)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	s, err := solver.FindSolution(ctx)
	if err != nil {
		log.Fatalf("%e", err)
	}

	fullRes := make([]float64, size)

	fmt.Println(s.RawVector().Data)

	err = comm.AllGatherF64(fullRes, s.RawVector().Data)
	if err != nil {
		log.Fatalf("can't record res data %e", err)
	}

	//fmt.Println(fullRes)

	log.Println("Program time: ", time.Since(startTime))
}

func initMpiMatrix(comm *mpi.Comm) (int, int, []float64) {
	start, end := pariter.GetMpiChunkParams(comm, size)

	if comm.Rank() == 0 {
		return start, end, generator.GenRandomVector(size)
	}

	return start, end, make([]float64, size)
}
