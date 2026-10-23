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
