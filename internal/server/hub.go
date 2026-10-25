package server

import (
	"log"
	"sync"
)

type Hub struct {
	mutex    sync.Mutex
	sessions map[*Session]struct{}
}

func NewHub() *Hub {
	return &Hub{
		sessions: make(map[*Session]struct{}),
	}
}

func (hub *Hub) Add(session *Session) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	hub.sessions[session] = struct{}{}
}

func (hub *Hub) Remove(session *Session) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	delete(hub.sessions, session)
}

func (hub *Hub) Count() int {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	return len(hub.sessions)
}

func (hub *Hub) Broadcast(message string) {
	hub.mutex.Lock()

	sessions := make([]*Session, 0, len(hub.sessions))
	for session := range hub.sessions {
		sessions = append(sessions, session)
	}

	hub.mutex.Unlock()
	for _, session := range sessions {
		err := session.Send(message)
		if err != nil {
			log.Printf("Unable to broadcast message: %s", err)
		}
	}
}

func (hub *Hub) snapshot() []*Session {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	sessions := make([]*Session, 0, len(hub.sessions))
	for session := range hub.sessions {
		sessions = append(sessions, session)
	}

	return sessions
}
