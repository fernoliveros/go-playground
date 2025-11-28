package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

	flip := false

	for {
		select {
		case <-ticker.C:
			sendHttpRequest(flip)
			flip = !flip
		case <-ctx.Done():
			fmt.Println("Context has been cancelled!! exiting")
			return
		}

	}
}

func sendHttpRequest(flip bool) {
	if flip {
		sendPostRequest()
		return
	}
	var url = "http://localhost:8080"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("ERROR GETTING FROM LOCALHOST 8080")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Got an error status code %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ERROR reading response body %v", err)
	}

	fmt.Printf("Body: \n%s\n", body)
}

func sendPostRequest() {
	var url = "http://localhost:8080/add"
	postData := "just a simple string"

	jsonBytes, err := json.Marshal(postData)
	if err != nil {
		fmt.Println("error marshaling data to json bytes")
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		fmt.Println("error sending post request")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Got an error status code %v", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ERROR reading response body %v", err)
	}

	fmt.Printf("Body: \n%s\n", body)
}
