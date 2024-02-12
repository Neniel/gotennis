package main

import (
	"context"
)

func main() {
	app := NewApp(context.Background())
	app.NewAPIServer().Run()
}
