package main

import (
	"context"

	"github.com/Neniel/gotennis/lib/app"
)

func main() {
	app := app.NewApp(context.Background())

	ms := &CustomerMicroservice{
		App: app,
	}

	ms.NewAPIServer().Run()
}
