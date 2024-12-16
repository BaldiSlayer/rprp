package tmeas

import (
	"fmt"
	"github.com/BaldiSlayer/rprp/lab5/internal/defaults"
	"time"
)

func MeasureFuncTime(logString string, f func()) {
	start := time.Now()

	f()

	sequentialDuration := time.Since(start)
	fmt.Printf("%s %v\n", logString, sequentialDuration/time.Duration(defaults.Iterations))
}
