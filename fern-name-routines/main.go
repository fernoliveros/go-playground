package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// type Fernando struct {
// 	firstName string
// 	lastName  string
// 	age       int
// }

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

// func someProcess(user *Fernando, key string, val string, strWg *sync.WaitGroup, errChan chan error) {
// 	fmt.Printf("Processing %s: %s\n", key, val)

// 	defer strWg.Done()

// 	time.Sleep(time.Second)

// 	switch key {
// 	case "firstName":
// 		user.firstName = val
// 	case "lastName":
// 		user.lastName = val
// 	case "age":
// 		ageint, err := strconv.Atoi(val)
// 		if err != nil {
// 			errChan <- fmt.Errorf("ERROR! Cannot cast string age %v to int", val)
// 			return
// 		}

// 		user.age = ageint
// 	}
// }
