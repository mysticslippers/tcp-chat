package server

import (
	"fmt"
	"log"
	"sync"
)

const (
	historyLimit    = 20
	defaultRoomName = "general"
)

type Hub struct {
	mutex    sync.Mutex
	sessions map[*Session]struct{}
	rooms    map[string]*Room
}

func NewHub() *Hub {
	return &Hub{
		sessions: make(map[*Session]struct{}),
		rooms: map[string]*Room{
			defaultRoomName: newRoom(defaultRoomName),
		},
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
		room := hub.rooms[defaultRoomName]
		room.sessions[session] = struct{}{}
		session.room = room
		break
	}

	history := make([]string, len(session.room.history))
	copy(history, session.room.history)

	return history
}

func (hub *Hub) Remove(session *Session) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	delete(hub.sessions, session)

	if session.room != nil {
		delete(session.room.sessions, session)
		session.room = nil
	}
}

func (hub *Hub) GetRoomNicknames(session *Session) []string {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	room := session.room
	if room == nil {
		return nil
	}

	nicknames := make([]string, 0, len(room.sessions))

	for participant := range room.sessions {
		nicknames = append(nicknames, participant.GetNickname())
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

func (hub *Hub) RestoreHistory(history map[string][]string) error {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	if len(hub.sessions) > 0 {
		return fmt.Errorf("cannot restore history with connected clients")
	}

	rooms := map[string]*Room{
		defaultRoomName: newRoom(defaultRoomName),
	}

	for name, messages := range history {
		if name == "" {
			return fmt.Errorf("room name must not be empty")
		}

		if len(messages) > historyLimit {
			messages = messages[len(messages)-historyLimit:]
		}

		room := newRoom(name)
		room.history = make([]string, len(messages))
		copy(room.history, messages)

		rooms[name] = room
	}

	hub.rooms = rooms

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

	hub.saveMessageLocked(hub.rooms[defaultRoomName], message)
}

func (hub *Hub) GetHistory() []string {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	room := hub.rooms[defaultRoomName]

	history := make([]string, len(room.history))
	copy(history, room.history)

	return history
}

func (hub *Hub) GetRoomHistory(session *Session) []string {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	room := session.room
	if room == nil {
		return nil
	}

	history := make([]string, len(room.history))
	copy(history, room.history)

	return history
}

func (hub *Hub) GetRoomName(session *Session) string {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	if session == nil || session.room == nil {
		return ""
	}

	return session.room.name
}

func (hub *Hub) GetRooms() map[string]int {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	rooms := make(map[string]int, len(hub.rooms))

	for name, room := range hub.rooms {
		rooms[name] = len(room.sessions)
	}

	return rooms
}

func (hub *Hub) ExportHistory() map[string][]string {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	history := make(map[string][]string, len(hub.rooms))

	for name, room := range hub.rooms {
		messages := make([]string, len(room.history))
		copy(messages, room.history)

		history[name] = messages
	}

	return history
}

func (hub *Hub) Publish(sender *Session, message string) {
	var failed []*Session

	hub.mutex.Lock()

	room := sender.room
	if room == nil {
		hub.mutex.Unlock()
		return
	}

	hub.saveMessageLocked(room, message)

	for recipient := range room.sessions {
		err := recipient.enqueue(message)
		if err != nil {
			failed = append(failed, recipient)
		}
	}

	hub.mutex.Unlock()

	for _, session := range failed {
		log.Printf(
			"Unable to queue message for %q; closing session",
			session.GetNickname(),
		)

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

func (hub *Hub) Join(session *Session, name string) error {
	hub.mutex.Lock()

	if name == "" {
		return fmt.Errorf("room name must not be empty")
	}

	_, exists := hub.sessions[session]

	if !exists {
		return fmt.Errorf("session is not registered")
	}

	room, exist := hub.rooms[name]
	if !exist {
		room = newRoom(name)
		hub.rooms[name] = room
	}

	if session.room == room {
		hub.mutex.Unlock()
		return nil
	}

	nickname := session.GetNickname()
	var failed []*Session

	if session.room != nil {
		old := session.room
		delete(old.sessions, session)

		message := fmt.Sprintf("[server] %s left %s\n", nickname, old.name)

		failed = hub.notifyRoomLocked(old, session, message)
	}

	room.sessions[session] = struct{}{}
	session.room = room

	message := fmt.Sprintf("[server] %s joined %s\n", nickname, room.name)

	failed = append(failed, hub.notifyRoomLocked(room, session, message)...)

	hub.mutex.Unlock()
	closeFailedSessions(failed)
	return nil
}

func (hub *Hub) isNicknameTaken(name string) bool {
	for session := range hub.sessions {
		if session.GetNickname() == name {
			return true
		}
	}

	return false
}

func (hub *Hub) saveMessageLocked(room *Room, message string) {
	room.history = append(room.history, message)

	if len(room.history) > historyLimit {
		room.history = room.history[1:]
	}
}
