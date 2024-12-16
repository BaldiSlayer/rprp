package main

import (
	"github.com/BaldiSlayer/rprp/lab5/internal/app"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(16)
	app.New().Run()
}
