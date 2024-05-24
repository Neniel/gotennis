package main

import (
	"context"

	"github.com/Neniel/gotennis/lib/app"
)

func main() {
	ms := &CacheMicroservice{
		App: app.NewApp(context.Background()),
	}
	ms.StartSync()
}
