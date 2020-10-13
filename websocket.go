package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocket _
type WebSocket struct {
	Conn   *websocket.Conn
	Out    chan []byte
	In     chan []byte
	Events map[string]EventHandler
}

// NewWebSocket _
func NewWebSocket(w http.ResponseWriter, r *http.Request) (*WebSocket, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("An error occured while upgrading the connection: %v", err)
		return nil, err
	}

	ws := &WebSocket{
		Conn:   conn,
		Out:    make(chan []byte),
		In:     make(chan []byte),
		Events: make(map[string]EventHandler),
	}

	go ws.Reader()
	go ws.Writer()

	return ws, nil
}

// Reader _
func (ws *WebSocket) Reader() {
	defer func() {
		ws.Conn.Close()
	}()

	for {
		_, message, err := ws.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("ws message error: %v", err)
			}
			break
		}

		event, err := NewEvent(message)
		if err != nil {
			log.Printf("error parsing message: %v", err)
			return
		}

		if action, ok := ws.Events[event.Name]; ok {
			action(event)
		}

	}

}

// Writer _
func (ws *WebSocket) Writer() {
	for {
		select {
		case message, ok := <-ws.Out:
			if !ok {
				ws.Conn.WriteMessage(websocket.CloseMessage, make([]byte, 0))
				return
			}
			w, err := ws.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			w.Close()
		}
	}
}

// On _
func (ws *WebSocket) On(event string, action EventHandler) *WebSocket {
	ws.Events[event] = action
	return ws
}
