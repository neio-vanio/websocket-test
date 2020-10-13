package main

import (
	"log"
	"net/http"
)

func main() {

	http.Handle("/", http.FileServer(http.Dir("./web")))

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws, err := NewWebSocket(w, r)
		if err != nil {
			log.Printf("Error creating websocket connection: %v", err)
		}

		ws.On("message", func(e *Event) {
			log.Printf("message received %s", e)
			ws.Out <- (&Event{
				Name: "received",
				Data: e.Data,
			}).Raw()
		})

	})

	log.Printf("websocket listening on port: %v", 3000)
	http.ListenAndServe(":3000", nil)
}
