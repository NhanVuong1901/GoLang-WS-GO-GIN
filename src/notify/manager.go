package notify

import (
	"sync"
	"ws/src/chat"
)

type NotifyHub struct {
	Clients    map[string]map[*chat.Client]bool
	Register   chan *chat.Client
	UnRegister chan *chat.Client
	Broadcast  chan *NotifyPayload
	mu         sync.Mutex
}

type NotifyPayload struct {
	UserID  string
	Message []byte
}

var NotifyWS = NewNotifyHub()

func NewNotifyHub() *NotifyHub {
	return &NotifyHub{
		Clients:    map[string]map[*chat.Client]bool{},
		Register:   make(chan *chat.Client),
		UnRegister: make(chan *chat.Client),
		Broadcast:  make(chan *NotifyPayload),
	}
}

func (h *NotifyHub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if h.Clients[client.UserID] == nil {
				h.Clients[client.UserID] = map[*chat.Client]bool{}
			}
			h.Clients[client.UserID][client] = true
			h.mu.Unlock()
		case client := <-h.UnRegister:
			h.mu.Lock()
			if _, ok := h.Clients[client.UserID]; ok {
				delete(h.Clients[client.UserID], client)
				close(client.Send)
			}
			h.mu.Unlock()
		case payload := <-h.Broadcast:
			h.mu.Lock()
			for client := range h.Clients[payload.UserID] {
				select {
				case client.Send <- payload.Message:
				default:
					close(client.Send)
					delete(h.Clients[client.UserID], client)
				}
			}
			h.mu.Unlock()
		}

	}
}