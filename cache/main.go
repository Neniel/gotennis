package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	//cacheSyncronizer := cacheSyncronizer{}
	//cacheSyncronizer.StartSync()
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg.Add(1)
	go sleep(ctx, &wg, 10)
	cancel()
	wg.Wait()
}

func sleep(ctx context.Context, wg *sync.WaitGroup, seconds int) {
	defer func() {
		fmt.Println("Exiting sleep")
	}()
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Printf("Sleeping for %d seconds\n", seconds)
			time.Sleep(time.Second * time.Duration(seconds))
		}
	}
}
