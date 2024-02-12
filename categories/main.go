package main

import (
	"context"

	"github.com/Neniel/gotennis/app"
)

func main() {
	app := app.NewApp(context.Background())

	ms := &CategoryMicroservice{
		App: app,
	}

	ms.NewAPIServer().Run()
}
