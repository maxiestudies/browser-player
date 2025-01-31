package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	// "time"

	"github.com/gorilla/websocket"
	"github.com/hypebeast/go-osc/osc"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	Action string `json:"action"`
	Src    string `json:"src,omitempty"`
}

func (h *Handler) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	log.Println("Client connected")

	var mu sync.Mutex

	// Handle incoming messages
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				return
			}
			log.Printf("Received: %s\n", message)
		}
	}()

	// Send commands to the client at specific intervals
	sendCommand := func(action string, src string) {
		mu.Lock()
		defer mu.Unlock()
		msg := Message{Action: action, Src: src}
		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			log.Println("JSON marshal error:", err)
			return
		}
		if err := conn.WriteMessage(websocket.TextMessage, jsonMsg); err != nil {
			log.Println("Write error:", err)
		}
	}

	// Wait for the connection to close
	for {
		msg := <-h.ch
		strMsg := msg.String()
		log.Printf("inside handler: %v", strMsg)
		command := msg.Arguments[0].(string)

		param := ""
		if msg.CountArguments() > 1 {
			param = msg.Arguments[1].(string)
		}

		sendCommand(command, param)
	}
}

type Handler struct {
	ch chan *osc.Message
}

func main() {

	oscChannel := make(chan *osc.Message)
	// osc
	addr := ":8765"
	d := osc.NewStandardDispatcher()
	d.AddMsgHandler("/osc/player1", func(msg *osc.Message) {
		// osc.PrintMessage(msg)
		oscChannel <- msg
		osc.PrintMessage(msg)
	})

	server := &osc.Server{
		Addr: addr,
		Dispatcher:d,
	}
	go server.ListenAndServe()

	handler := Handler{
		ch: oscChannel,
	}

	// Serve static files from the "public" directory
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	// WebSocket endpoint
	http.HandleFunc("/ws", handler.handleWebSocket)

	// Start the server
	log.Println("Server is listening on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}

