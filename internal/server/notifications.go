package server

import "log"

func closeFailedSessions(sessions []*Session) {
	for _, session := range sessions {
		log.Printf("Unable to queue message for %q; closing session", session.GetNickname())

		session.close()
	}
}
