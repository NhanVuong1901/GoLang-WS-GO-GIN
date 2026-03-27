package room

import (
	"sync"
)

var RoomMembers = NewPresenceTracker()

type PresenceTracker struct {
	mu     sync.Mutex
	online map[string]map[string]bool // roomID -> map[userID]bool
}

func NewPresenceTracker() *PresenceTracker {
	return &PresenceTracker{
		online: make(map[string]map[string]bool),
	}
}
