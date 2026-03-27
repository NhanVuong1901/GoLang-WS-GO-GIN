package chat

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn   *websocket.Conn
	UserID string
	RoomID string
	Send   chan []byte
}

type Hub struct {
	Clients    map[string]map[*Client]bool
	Register   chan *Client
	UnRegister chan *Client
	Broadcast  chan *MessagePayload
	mu         sync.Mutex
}

type MessagePayload struct {
	RoomID  string
	Message []byte
}

var WS = NewHub()

func NewHub() *Hub {
	return &Hub{
		Clients:    map[string]map[*Client]bool{},
		Register:   make(chan *Client),
		UnRegister: make(chan *Client),
		Broadcast:  make(chan *MessagePayload),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if h.Clients[client.RoomID] == nil {
				h.Clients[client.RoomID] = map[*Client]bool{}
			}
			h.Clients[client.RoomID][client] = true
			h.mu.Unlock()
		case client := <-h.UnRegister:
			h.mu.Lock()
			if _, ok := h.Clients[client.RoomID]; ok {
				delete(h.Clients[client.RoomID], client)
				close(client.Send)
			}
			h.mu.Unlock()
		case payload := <-h.Broadcast:
			h.mu.Lock()
			for client := range h.Clients[payload.RoomID] {
				select {
				case client.Send <- payload.Message:
				default:
					close(client.Send)
					delete(h.Clients[client.RoomID], client)
				}
			}
			h.mu.Unlock()
		}
	}
}
