package server

type Room struct {
	name     string
	history  []string
	sessions map[string]*Session
}

func newRoom(name string) *Room {
	return &Room{
		name:     name,
		sessions: make(map[string]*Session),
	}
}
