package server

type Room struct {
	name     string
	history  []string
	sessions map[*Session]struct{}
}

func newRoom(name string) *Room {
	return &Room{
		name:     name,
		sessions: make(map[*Session]struct{}),
	}
}
