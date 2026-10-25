package server

import (
	"fmt"
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

	for number := 1; ; number++ {
		name := fmt.Sprintf("guest-%d", number)

		if hub.isNicknameTaken(name) {
			continue
		}

		session.SetNickname(name)
		hub.sessions[session] = struct{}{}
		return
	}
}

func (hub *Hub) Remove(session *Session) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	delete(hub.sessions, session)
}

func (hub *Hub) GetNicknames() []string {
	sessions := hub.snapshot()
	nicknames := make([]string, 0, len(sessions))

	for _, session := range sessions {
		nicknames = append(nicknames, session.nickname)
	}

	return nicknames
}

func (hub *Hub) Count() int {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	return len(hub.sessions)
}

func (hub *Hub) Broadcast(message string) {
	sessions := hub.snapshot()

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

func (hub *Hub) isNicknameTaken(name string) bool {
	for session := range hub.sessions {
		if session.GetNickname() == name {
			return true
		}
	}

	return false
}
