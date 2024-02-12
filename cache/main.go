package main

import (
	"context"

	"github.com/Neniel/gotennis/app"
)

func main() {
	ms := &CacheMicroservice{
		App: app.NewApp(context.Background()),
	}
	ms.StartSync()
}
