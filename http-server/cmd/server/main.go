package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("here we go")

	server := &http.Server{Addr: ":8080"}
	http.HandleFunc("/", defaultHanler)
	http.HandleFunc("/add", addHanler)

	go httpServer(server)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	<-stop

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Failed to shut down http server gracefully")
	}
	fmt.Println("Gracefully shutting down the http server")

}

func defaultHanler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Default Handler")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	fmt.Fprintf(w, "Thanks for sending us a default request")
}
func addHanler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Add Handler")

	if r.Method != "POST" {
		http.Error(w, "we only allow posts", http.StatusMethodNotAllowed)
	}

	jsonBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("ERROR reading response body %v", err)
	}

	var body string
	if err := json.Unmarshal(jsonBytes, &body); err != nil {
		log.Fatalf("ERROR unmarshalling response body %v", err)

	}
	fmt.Printf("Received post with %v\n", body)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	fmt.Fprintf(w, "Thanks for calling the add hanlder, you sent a post!")
}

func httpServer(server *http.Server) {

	fmt.Println("HTTP Server listening on 8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to listen and serve on 8080 %v", err)
	}
}
