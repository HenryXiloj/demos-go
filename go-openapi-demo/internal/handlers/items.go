package handlers

import (
	"encoding/json"
	"net/http"
	"sync"
	"sync/atomic"

	"go-openapi-demo/api"
)

type ItemsService struct {
	mu        sync.RWMutex
	items     map[int64]api.Item
	currentID int64
}

func NewItemsService() *ItemsService {
	s := &ItemsService{
		items:     map[int64]api.Item{},
		currentID: 2,
	}
	s.items[1] = api.Item{Id: 1, Name: "Item 1"}
	s.items[2] = api.Item{Id: 2, Name: "Item 2"}
	return s
}

func (s *ItemsService) GetItems(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	out := make([]api.Item, 0, len(s.items))
	for _, it := range s.items {
		out = append(out, it)
	}
	s.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(out)
}

func (s *ItemsService) GetItemById(w http.ResponseWriter, r *http.Request, id int64) {
	s.mu.RLock()
	it, ok := s.items[id]
	s.mu.RUnlock()
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(it)
}

func (s *ItemsService) CreateItem(w http.ResponseWriter, r *http.Request) {
	var req api.ItemCreate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	id := atomic.AddInt64(&s.currentID, 1)
	item := api.Item{Id: id, Name: req.Name}

	s.mu.Lock()
	s.items[id] = item
	s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(item)
}

func (s *ItemsService) UpdateItem(w http.ResponseWriter, r *http.Request, id int64) {
	var req api.ItemUpdate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[id]; !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	s.items[id] = api.Item{Id: id, Name: req.Name}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(s.items[id])
}

func (s *ItemsService) DeleteItem(w http.ResponseWriter, r *http.Request, id int64) {
	s.mu.Lock()
	_, ok := s.items[id]
	if ok {
		delete(s.items, id)
	}
	s.mu.Unlock()

	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
