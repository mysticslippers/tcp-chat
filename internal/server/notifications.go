package server

import "log"

func closeFailedSessions(sessions []*Session) {
	for _, session := range sessions {
		log.Printf("Unable to queue message for %q; closing session", session.GetNickname())

		session.close()
	}
}

func (hub *Hub) notifyRoomLocked(room *Room, excluded *Session, message string) []*Session {
	var failed []*Session

	for recipient := range room.sessions {
		if recipient == excluded {
			continue
		}

		err := recipient.enqueue(message)
		if err != nil {
			failed = append(failed, recipient)
		}
	}

	return failed
}
