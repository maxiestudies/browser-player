package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

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

	log.Printf("Client %s connected", conn.RemoteAddr().String())

	var mu sync.Mutex
	_, player, err := conn.ReadMessage()
	if err != nil {
		log.Println("Read error:", err)
	}
	log.Printf("Player %s is ready", string(player))

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
		msg := <-h.playerChannels[string(player)]
		strMsg := msg.String()
		log.Printf("inside handler: %v", strMsg)
		command := msg.Arguments[1].(string)

		param := ""
		if msg.CountArguments() > 2 {
			param = msg.Arguments[2].(string)
		}

		sendCommand(command, param)
	}
}

type Handler struct {
	ch chan *osc.Message
	playerChannels map[string]chan *osc.Message
}

func main() {

	oscChannel := make(chan *osc.Message)
	// playerChannels := make(map[string]chan *osc.Message)
	playerChannels := make(map[string]chan *osc.Message)
	playerChannels["player1"] = make(chan *osc.Message)
	playerChannels["player2"] = make(chan *osc.Message)

	// osc
	addr := ":8765"
	d := osc.NewStandardDispatcher()
	d.AddMsgHandler("/osc", func(msg *osc.Message) {
		osc.PrintMessage(msg)
		player := msg.Arguments[0].(string)
		playerChannels[player] <- msg
		// osc.PrintMessage(msg)
		// oscChannel <- msg
	})

	server := &osc.Server{
		Addr: addr,
		Dispatcher:d,
	}
	go server.ListenAndServe()

	handler := Handler{
		ch: oscChannel,
		playerChannels: playerChannels,
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

