package main

import (
	"log"
	"runtime"
	"time"

	"github.com/BaldiSlayer/rprp/lab1/internal/experiment"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

const (
	maxValues         = 200
	some          int = 1650
	maxGoroutines     = 7
)

func durationToFloat64(s time.Duration) float64 {
	return s.Seconds()
}

func typeFind(i int) string {
	if i == 0 {
		return "Normal"
	}

	if i == 1 {
		return "Cols"
	}

	if i == 2 {
		return "Rows"
	}

	if i == 3 {
		return "Blocks"
	}

	return ""
}

func main() {
	runtime.GOMAXPROCS(2)

	log.Print(runtime.GOMAXPROCS(0))

	exp := experiment.New(2*some, 1*some, maxValues, 1*some, 2*some, maxValues)
	log.Print("Начало эксперимента")

	lines := make([]plotter.XYs, 4)
	for g := 1; g < maxGoroutines; g++ {
		log.Printf("Количество горутин: %d", g)

		results, err := exp.Run(g)
		if err != nil {
			log.Fatalf("error while running experiment: %v", err)
		}

		lines[0] = append(lines[0], plotter.XY{X: float64(g), Y: durationToFloat64(results.Normal)})
		lines[1] = append(lines[1], plotter.XY{X: float64(g), Y: durationToFloat64(results.Cols)})
		lines[2] = append(lines[2], plotter.XY{X: float64(g), Y: durationToFloat64(results.Rows)})
		lines[3] = append(lines[3], plotter.XY{X: float64(g), Y: durationToFloat64(results.Blocks)})
	}

	p := plot.New()
	p.Title.Text = "Результаты"
	p.X.Label.Text = "Количество потоков"
	p.Y.Label.Text = "Время в миллисекундах"

	for i, points := range lines {
		line, err := plotter.NewLine(points)
		if err != nil {
			log.Fatalf("%v", err)
		}

		colorIndex := i % len(plotutil.SoftColors)
		line.Color = plotutil.SoftColors[colorIndex]

		p.Legend.Add(typeFind(i), line)
		p.Add(line)
	}

	p.Legend.Top = true

	// Сохраняем график в файл.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "plot.png"); err != nil {
		log.Fatalf("ошибка при сохранении в файл: %v", err)
	}
}
