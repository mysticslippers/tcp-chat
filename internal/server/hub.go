package server

import "sync"

type Hub struct {
	mutex    sync.Mutex
	sessions map[*Session]struct{}
}
