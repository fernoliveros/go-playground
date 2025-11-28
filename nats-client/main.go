package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	nats "github.com/nats-io/nats.go"
)

type Messages struct {
	Messages []string
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Required 3 args: sub chan, pub chan and http port")
		os.Exit(1)
	}

	natSubChannel := os.Args[1]
	natPubChannel := os.Args[2]
	httpPort := os.Args[3]

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		fmt.Printf("Error connecting to nats server: %v", err)
		return
	}

	// subscribeToNatChannel(nc, natSubChannel)
	registerPublishToNatChannelHandler(nc, natPubChannel)
	setupMessagesStream(nc, natSubChannel)

	// var wg sync.WaitGroup
	// wg.Add(1)
	startHtmxServer(httpPort)
	// wg.Wait()

}

func setupMessagesStream(nc *nats.Conn, natSubChannel string) {
	http.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)

		fmt.Println("The /messages handler has been invoked!")

		ch := make(chan *nats.Msg, 64)
		_, err := nc.ChanSubscribe(natSubChannel, ch)
		if err != nil {
			fmt.Printf("Error subscribing to the natsubChannel: %v", err.Error())
		}

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Error initializing flusher, streaming may be unsupported", http.StatusInternalServerError)
			return
		}

		ctx := r.Context()

		msgSlice := make([]string, 0, 10)
		messagesForTmpl := &Messages{
			Messages: msgSlice,
		}

		for {
			select {
			case <-ctx.Done():
				fmt.Println("Client disconnected")
			case msg := <-ch:

				fmt.Sprintf("msg not used but received: %s", string(msg.Data))

				messagesForTmpl.Messages = append(messagesForTmpl.Messages, string(msg.Data))
				messagesTemplate, err := buildMessagesTemplate(messagesForTmpl)
				if err != nil {
					fmt.Printf("Error building messages template: %v", err.Error())
				}

				ssePayload := fmt.Sprintf("event:message\ndata: %s\n\n", messagesTemplate)

				fmt.Printf("About to write to the messages stream:\n %s", ssePayload)
				_, err = w.Write([]byte(ssePayload))
				if err != nil {
					fmt.Printf("Error writting data to the messages stream: %v", err.Error())
				}
				flusher.Flush()
			}
		}

		// nc.Subscribe(natSubChannel, func(m *nats.Msg) {
		// 	fmt.Printf(":: %s\n", string(m.Data))

		// })
	})
}

func startHtmxServer(port string) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("html/index.html")
		if err != nil {
			fmt.Printf("Error parsing index.html %v", err.Error())
			http.Error(w, "Error parsing index.html", http.StatusInternalServerError)
			return
		}

		data := struct{ Message string }{
			Message: "Hello World!",
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			fmt.Printf("Error executing data struct into index.html %v \n", err.Error())
			http.Error(w, "Error executing data struct into index.html", http.StatusInternalServerError)
			return
		}
	})
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

// func subscribeToNatChannel(nc *nats.Conn, natSubChannel string) {
// 	nc.Subscribe(natSubChannel, func(m *nats.Msg) {
// 		fmt.Printf(":: %s\n", string(m.Data))
// 	})
// }

func registerPublishToNatChannelHandler(nc *nats.Conn, natPubChannel string) {
	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		nc.Publish(natPubChannel, []byte("Hello there!"))
	})

}

func buildMessagesTemplate(messagesForTmpl *Messages) ([]byte, error) {
	tmpl, err := template.ParseFiles("html/messages.html")
	if err != nil {
		fmt.Printf("Error parsing messages.html %v", err.Error())
		return nil, fmt.Errorf("Error parsing messages.html %w", err)
	}

	var buf strings.Builder

	err = tmpl.Execute(&buf, messagesForTmpl)
	if err != nil {
		fmt.Printf("Error executing data struct into messages.html %v \n", err.Error())
		return nil, fmt.Errorf("Error executing data struct into messages.html %w", err)
	}

	return []byte(buf.String()), nil
}
