package server

import (
	"fmt"
	"log"
	"sync"
)

const historyLimit = 20

type Hub struct {
	mutex    sync.Mutex
	history  []string
	sessions map[*Session]struct{}
}

func NewHub() *Hub {
	return &Hub{
		sessions: make(map[*Session]struct{}),
	}
}

func (hub *Hub) Add(session *Session) []string {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	for number := 1; ; number++ {
		name := fmt.Sprintf("guest-%d", number)

		if hub.isNicknameTaken(name) {
			continue
		}

		session.setNickname(name)
		hub.sessions[session] = struct{}{}
		break
	}

	history := make([]string, len(hub.history))
	copy(history, hub.history)

	return history
}

func (hub *Hub) Remove(session *Session) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	delete(hub.sessions, session)
}

func (hub *Hub) GetNicknames() []string {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	nicknames := make([]string, 0, len(hub.sessions))

	for session := range hub.sessions {
		nicknames = append(nicknames, session.nickname)
	}

	return nicknames
}

func (hub *Hub) Rename(session *Session, name string) error {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	_, exist := hub.sessions[session]
	if !exist {
		return fmt.Errorf("session is not registered")
	}

	if session.GetNickname() == name {
		return nil
	}

	if hub.isNicknameTaken(name) {
		return fmt.Errorf("nickname %q is already taken", name)
	}

	session.setNickname(name)
	return nil
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

func (hub *Hub) FindByNickname(nickname string) *Session {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	for session := range hub.sessions {
		if session.GetNickname() == nickname {
			return session
		}
	}

	return nil
}

func (hub *Hub) SaveMessage(message string) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	hub.saveMessageLocked(message)
}

func (hub *Hub) GetHistory() []string {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	history := make([]string, len(hub.history))
	copy(history, hub.history)

	return history
}

func (hub *Hub) Publish(message string) {
	var failed []*Session

	hub.mutex.Lock()
	hub.saveMessageLocked(message)

	for session := range hub.sessions {
		err := session.enqueue(message)
		if err != nil {
			failed = append(failed, session)
		}
	}

	hub.mutex.Unlock()
	for _, session := range failed {
		log.Printf("Unable to queue message for %q; closing session", session.GetNickname())

		session.close()
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

func (hub *Hub) saveMessageLocked(message string) {
	hub.history = append(hub.history, message)

	if len(hub.history) > historyLimit {
		hub.history = hub.history[1:]
	}
}
