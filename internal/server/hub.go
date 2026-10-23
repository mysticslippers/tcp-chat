package server

import "sync"

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
