package bot

import "sync"

type Store struct {
	mu     sync.RWMutex
	bots   map[string]Bot
	latest map[string]Message
	status map[string]Status
}

func NewStore() *Store {
	return &Store{
		bots:   make(map[string]Bot),
		latest: make(map[string]Message),
		status: make(map[string]Status),
	}
}

func (s *Store) Exists(id string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.bots[id]
	return ok
}

func (s *Store) Add(b Bot) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.bots[b.Id] = b
}

func (s *Store) SetLatest(targetId string, msg Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.latest[targetId] = msg
}

// GetLatest returns the most recent instruction, if any.
func (s *Store) GetLatest(id string) (Message, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	msg, ok := s.latest[id]
	return msg, ok
}

func (s *Store) SetStatus(id string, st Status) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status[id] = st
}

func (s *Store) GetStatus(id string) (Status, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	st, ok := s.status[id]
	return st, ok
}
