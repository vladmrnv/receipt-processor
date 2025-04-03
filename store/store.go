package store

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/receipt-processor/models"
)

var ErrReceiptNotFound = errors.New("receipt not found")

type Store struct {
	receipts map[string]models.Receipt
	mu       sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		receipts: make(map[string]models.Receipt),
	}
}

func (s *Store) SaveReceipt(receipt models.Receipt) string {
	id := uuid.New().String()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.receipts[id] = receipt
	return id
}

func (s *Store) GetReceipt(id string) (models.Receipt, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	receipt, ok := s.receipts[id]
	if !ok {
		return models.Receipt{}, ErrReceiptNotFound
	}
	return receipt, nil
}
