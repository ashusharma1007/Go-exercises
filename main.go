//go:build !learning
// +build !learning

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for local development
	},
}

// Message types
type Message struct {
	Type      string    `json:"type"`
	Operation Operation `json:"operation,omitempty"`
	Content   string    `json:"content,omitempty"`
	Version   int       `json:"version,omitempty"`
	ClientID  string    `json:"clientId,omitempty"`
}

// Client represents a connected WebSocket client
type Client struct {
	ID   string
	Conn *websocket.Conn
	Send chan Message
}

// Hub maintains the document and manages clients
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	document   *Document
	mutex      sync.RWMutex
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		document:   &Document{Content: "", Version: 0},
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.mutex.Unlock()

			// Send current document state to new client
			client.Send <- Message{
				Type:    "init",
				Content: h.document.Content,
				Version: h.document.Version,
			}
			log.Printf("Client %s connected. Total clients: %d", client.ID, len(h.clients))

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
				log.Printf("Client %s disconnected. Total clients: %d", client.ID, len(h.clients))
			}
			h.mutex.Unlock()

		case message := <-h.broadcast:
			if message.Type == "operation" {
				h.mutex.Lock()

				// Apply operation to document
				err := h.document.Apply(message.Operation)
				if err != nil {
					log.Printf("Error applying operation: %v", err)
					h.mutex.Unlock()
					continue
				}

				log.Printf("Applied %s at pos %d. New content: '%s' (v%d)",
					message.Operation.Type, message.Operation.Pos,
					h.document.Content, h.document.Version)

				h.mutex.Unlock()
			}

			// Broadcast to all clients except sender
			h.mutex.RLock()
			for client := range h.clients {
				if client.ID != message.ClientID {
					select {
					case client.Send <- message:
					default:
						close(client.Send)
						delete(h.clients, client)
					}
				}
			}
			h.mutex.RUnlock()
		}
	}
}

func (c *Client) readPump(hub *Hub) {
	defer func() {
		hub.unregister <- c
		c.Conn.Close()
	}()

	for {
		var msg Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		msg.ClientID = c.ID
		hub.broadcast <- msg
	}
}

func (c *Client) writePump() {
	defer c.Conn.Close()

	for message := range c.Send {
		err := c.Conn.WriteJSON(message)
		if err != nil {
			log.Printf("Write error: %v", err)
			break
		}
	}
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	clientID := r.URL.Query().Get("clientId")
	if clientID == "" {
		clientID = fmt.Sprintf("client-%d", len(hub.clients)+1)
	}

	client := &Client{
		ID:   clientID,
		Conn: conn,
		Send: make(chan Message, 256),
	}

	hub.register <- client

	go client.writePump()
	go client.readPump(hub)
}

func main() {
	hub := newHub()
	go hub.run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	port = ":" + port
	log.Printf("Server starting on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
