// Package todo provides an in-memory to-do store and its HTTP handlers.
package todo

import (
	"errors"
	"sort"
	"sync"
	"time"
)

// ErrNotFound is returned when a to-do with the given ID does not exist.
var ErrNotFound = errors.New("todo not found")

// Todo is a single to-do item.
type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"createdAt"`
}

// Store is a concurrency-safe in-memory collection of to-dos.
type Store struct {
	mu     sync.RWMutex
	nextID int
	items  map[int]Todo
}

// NewStore returns an empty store.
func NewStore() *Store {
	return &Store{nextID: 1, items: make(map[int]Todo)}
}

// List returns all to-dos ordered by ID.
func (s *Store) List() []Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]Todo, 0, len(s.items))
	for _, t := range s.items {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// Get returns the to-do with the given ID.
func (s *Store) Get(id int) (Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, ok := s.items[id]
	if !ok {
		return Todo{}, ErrNotFound
	}
	return t, nil
}

// Create adds a new to-do and returns it.
func (s *Store) Create(title string) Todo {
	s.mu.Lock()
	defer s.mu.Unlock()

	t := Todo{ID: s.nextID, Title: title, CreatedAt: time.Now().UTC()}
	s.items[t.ID] = t
	s.nextID++
	return t
}

// Update changes the title and/or done flag of a to-do. Nil fields are left unchanged.
func (s *Store) Update(id int, title *string, done *bool) (Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, ok := s.items[id]
	if !ok {
		return Todo{}, ErrNotFound
	}
	if title != nil {
		t.Title = *title
	}
	if done != nil {
		t.Done = *done
	}
	s.items[id] = t
	return t, nil
}

// Delete removes a to-do.
func (s *Store) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[id]; !ok {
		return ErrNotFound
	}
	delete(s.items, id)
	return nil
}
