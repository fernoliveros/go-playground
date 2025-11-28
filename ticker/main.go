package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("Hello World!")

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var wg sync.WaitGroup
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())

	go foreverProcess(ticker, &wg, ctx)

	go func() {
		time.Sleep(time.Second * 10)
		cancel()
	}()

	wg.Wait()

}

func foreverProcess(ticker *time.Ticker, wg *sync.WaitGroup, ctx context.Context) {
	fmt.Println("foreverProcess started")
	defer wg.Done()

	secsPassed := 1

	for {
		select {
		case <-ticker.C:
			fmt.Println(secsPassed, "seconds passed")
			secsPassed += 1
		case <-ctx.Done():
			fmt.Println("Context has been cancelled!! exiting")
			return
		}

	}
}
