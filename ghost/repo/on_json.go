package repo

import (
	"io"
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type OnJson[T any] struct {
	file *os.File
	db   jsonDB[T]
	mu   sync.RWMutex
	setId  func(*T, int64)
}
type jsonDB[T any] struct {
	NextId int64       `json:"nextId"`
	Items  map[int64]T `json:"items"`
}

func NewOnJson[T any](path string, setId func(*T, int64)) (*OnJson[T], error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, errors.New("failed to open file")
	}

	newRepo := &OnJson[T]{
		file: file,
		setId: setId,
	}

	err = json.NewDecoder(file).Decode(&newRepo.db)

	if err != nil && !errors.Is(err, io.EOF) {
		file.Close()
		return nil, err
	}

	if newRepo.db.Items == nil {
		newRepo.db.Items = make(map[int64]T)
	}

	return newRepo, nil
}
// TODO: this save logic is very **danger**
func (s *OnJson[T]) Save() error {
	_, err := s.file.Seek(0, 0)
	if err != nil {
		return err
	}

	err = s.file.Truncate(0)
	if err != nil {
		return err
	}

	return json.NewEncoder(s.file).Encode(s.db)
}
func (s *OnJson[T]) Close() {
	s.Save()
	s.file.Close()
}
func (s *OnJson[T]) Create(t T) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.db.NextId++
	s.setId(&t, s.db.NextId)
	s.db.Items[s.db.NextId] = t

	err := s.Save()
	if err != nil {
		return 0, err
	}

	return s.db.NextId, nil
}
func (s *OnJson[T]) List() ([]T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ts := make([]T, 0)
	for _, value := range s.db.Items {
		ts = append(ts, value)
	}
	return ts, nil
}
func (s *OnJson[T]) Find(predicate func(T) bool) (T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, value := range s.db.Items {
		if predicate(value) {
			return value, nil
		}
	}

	var zero T
	return zero, errors.New("item not found")
}
func (s *OnJson[T]) FindMany(predicate func(T) bool) ([]T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []T
	for _, value := range s.db.Items {
		if predicate(value) {
			result = append(result, value)
		}
	}

	return result, nil
}
func (s *OnJson[T]) Get(id int64) (T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.db.Items[id]
	if !ok {
		var zero T
		return zero, errors.New("item not found")
	}
	return t, nil
}
func (s *OnJson[T]) Update(id int64, value T) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.db.Items[id]; !ok {
		return errors.New("item not found")
	}
	s.setId(&value, id)
	s.db.Items[id] = value

	err := s.Save()
	if err != nil {
		return err
	}

	return nil
}
func (s *OnJson[T]) Delete(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.db.Items[id]
	if !ok {
		return errors.New("item not found")
	}
	delete(s.db.Items, id)

	err := s.Save()
	if err != nil {
		return err
	}

	return nil
}
