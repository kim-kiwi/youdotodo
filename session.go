package main

import (
	"errors"
	"github.com/google/uuid"
	"sync"
)

type SessionManager struct {
	mu       sync.RWMutex
	sessions map[string]any
}

func NewSessionManager() *SessionManager {
	return &SessionManager{sessions: make(map[string]any)}
}

func (s *SessionManager) New() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := uuid.New().String()
	s.sessions[id] = nil
	return id
}

func (s *SessionManager) Set(id string, value any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.sessions[id]
	if !ok {
		return errors.New("Session not found")
	}
	s.sessions[id] = value
	return nil
}

func (s *SessionManager) Get(id string) (any, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.sessions[id]
	if !ok {
		return nil, errors.New("Session not found")
	}
	return value, nil
}

func (s *SessionManager) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.sessions[id]
	if !ok {
		return errors.New("Session not found")
	}
	delete(s.sessions, id)
	return nil
}
