package server

type Room struct {
	name     string
	history  []string
	sessions map[string]*Session
}
