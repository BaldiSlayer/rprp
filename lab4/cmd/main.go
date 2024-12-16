package main

import (
	"context"
	app2 "github.com/BaldiSlayer/rprp/lab4/internal/app"
	"log"
	"time"
)

var (
	ctx = context.Background()
)

func main() {
	app := app2.New()

	ctx2, cancel := context.WithTimeout(ctx, 40*time.Second)
	defer cancel()

	err := app.Run(ctx2, 5)
	if err != nil {
		log.Printf("some philosophers are hungry :(\nreason:%v", err)
	}

	log.Print("end of philosophers lives")
}
