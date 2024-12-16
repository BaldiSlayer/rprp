package pariter

import (
	"context"
	"errors"
	"math"

	"github.com/emer/empi/mpi"
	"gonum.org/v1/gonum/mat"
)

type mpiObj struct {
	comm  *mpi.Comm
	rank  int
	start int
	end   int
}

type matrixChunk struct {
	baseMatrixChunk *mat.Dense
	freeConst       *mat.VecDense
	size            int
}

type SolverWithVecSeparation struct {
	mo  *mpiObj
	mc  *matrixChunk
	tau float64
	eps float64
}

func NewSolverWithVecSeparation(
	comm *mpi.Comm, baseMatrixChunk *mat.Dense, freeConst *mat.VecDense, size int, eps float64,
) *SolverWithVecSeparation {
	return &SolverWithVecSeparation{
		mo:  newMpiObj(comm, size),
		mc:  newMatrixChunk(baseMatrixChunk, freeConst, size),
		eps: eps,
		tau: float64(1) / float64(size),
	}
}

func newMpiObj(comm *mpi.Comm, size int) *mpiObj {
	rowStart, rowEnd := GetMpiChunkParams(comm, size)

	return &mpiObj{
		rank:  comm.Rank(),
		start: rowStart,
		end:   rowEnd,
		comm:  comm,
	}
}

func GetMpiChunkParams(comm *mpi.Comm, size int) (int, int) {
	chunkSize := size / comm.Size()

	return comm.Rank() * chunkSize, (comm.Rank() + 1) * chunkSize
}

func newMatrixChunk(baseMatrixChunk *mat.Dense, freeConst *mat.VecDense, size int) *matrixChunk {
	return &matrixChunk{
		baseMatrixChunk: baseMatrixChunk,
		freeConst:       freeConst,
		size:            size,
	}
}

func (s *SolverWithVecSeparation) FindSolution(ctx context.Context) (*mat.VecDense, error) {
	res := mat.NewVecDense(s.mo.end-s.mo.start, nil)
	buff := make([]float64, s.mc.size)
	iteration := 0

	err := s.iterateLoop(ctx, res, &buff, &iteration)

	return res, err
}

func (s *SolverWithVecSeparation) iterateLoop(
	ctx context.Context,
	res *mat.VecDense,
	buff *[]float64,
	iteration *int,
) error {
	metric := math.MaxFloat64
	prev := metric

	for metric > s.eps {
		select {
		case <-ctx.Done():
			return errors.New("context deadline")
		default:
		}

		*iteration++
		chunk := s.iterateSolution(res).RawVector().Data

		err := s.mo.comm.AllGatherF64(*buff, chunk)
		if err != nil {
			return err
		}

		tmp := mat.NewVecDense(s.mc.size, *buff)

		res.SubVec(res, tmp.SliceVec(s.mo.start, s.mo.end))

		prev = metric

		metric, err = s.calcMetric(res)
		if err != nil {
			return err
		}

		if prev < metric {
			s.tau *= -1
		}
	}

	return nil
}

func (s *SolverWithVecSeparation) iterateSolution(chunk *mat.VecDense) *mat.VecDense {
	buff := make([]float64, s.mc.size)
	_ = s.mo.comm.AllGatherF64(buff, chunk.RawVector().Data)

	tmp := mat.NewVecDense(s.mo.end-s.mo.start, nil)

	tmp.MulVec(s.mc.baseMatrixChunk, mat.NewVecDense(s.mc.size, buff))
	tmp.SubVec(tmp, s.mc.freeConst)
	tmp.ScaleVec(s.tau, tmp)

	return tmp
}

type fraction struct {
	numerator   float64
	denominator float64
}

func (s *SolverWithVecSeparation) calcMetric(chunk *mat.VecDense) (float64, error) {
	buff := make([]float64, s.mc.size)

	err := s.mo.comm.AllGatherF64(buff, chunk.RawVector().Data)
	if err != nil {
		return 0, err
	}

	tmp := mat.NewVecDense(s.mo.end-s.mo.start, nil)

	tmp.MulVec(s.mc.baseMatrixChunk, mat.NewVecDense(s.mc.size, buff))
	tmp.SubVec(tmp, s.mc.freeConst)

	frac := getChunkMetricFraction(s, tmp)

	sumFrac := make([]float64, 2)

	err = s.mo.comm.AllReduceF64(mpi.OpSum, sumFrac, []float64{frac.numerator, frac.denominator})
	if err != nil {
		return 0, err
	}

	return getMetricValue(sumFrac), nil
}

func getChunkMetricFraction(s *SolverWithVecSeparation, tmp *mat.VecDense) fraction {
	numerator, denominator := float64(0), float64(0)

	for _, v := range tmp.RawVector().Data {
		numerator += v * v
	}

	for _, v := range s.mc.freeConst.RawVector().Data {
		denominator += v * v
	}

	return fraction{
		numerator:   numerator,
		denominator: denominator,
	}
}

func getMetricValue(frac []float64) float64 {
	return math.Sqrt(frac[0]) / math.Sqrt(frac[1])
}
