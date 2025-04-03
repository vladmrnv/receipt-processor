package store

import (
	"fmt"
	"sync"
	"testing"

	"github.com/receipt-processor/models"
)

func TestStoreOperations(t *testing.T) {
	store := NewStore()

	receipt := models.Receipt{
		Retailer:     "TestStore",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []models.Item{
			{ShortDescription: "Test Item", Price: "10.00"},
		},
		Total: "10.00",
	}

	t.Run("SaveReceipt", func(t *testing.T) {
		id := store.SaveReceipt(receipt)
		if id == "" {
			t.Errorf("Expected non-empty ID, got empty string")
		}

		savedReceipt, err := store.GetReceipt(id)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if savedReceipt.Retailer != receipt.Retailer {
			t.Errorf("Expected retailer %s, got %s", receipt.Retailer, savedReceipt.Retailer)
		}

		if savedReceipt.Total != receipt.Total {
			t.Errorf("Expected total %s, got %s", receipt.Total, savedReceipt.Total)
		}

		if len(savedReceipt.Items) != len(receipt.Items) {
			t.Errorf("Expected %d items, got %d", len(receipt.Items), len(savedReceipt.Items))
		}
	})

	t.Run("GetNonExistentReceipt", func(t *testing.T) {
		_, err := store.GetReceipt("non-existent-id")
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("ConcurrentOperations", func(t *testing.T) {
		const numGoroutines = 10
		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		errorCh := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer wg.Done()
				id := store.SaveReceipt(receipt)
				if id == "" {
					errorCh <- fmt.Errorf("got empty ID during concurrent operation")
					return
				}

				_, err := store.GetReceipt(id)
				if err != nil {
					errorCh <- fmt.Errorf("unable to retrieve receipt during concurrent operation: %v", err)
					return
				}
			}()
		}

		wg.Wait()
		close(errorCh)

		for err := range errorCh {
			t.Errorf("Concurrent operation error: %v", err)
		}
	})
}

func TestStoreIDUniqueness(t *testing.T) {
	store := NewStore()
	receipt := models.Receipt{
		Retailer:     "TestStore",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Items: []models.Item{
			{ShortDescription: "Test Item", Price: "10.00"},
		},
		Total: "10.00",
	}

	const idCount = 100
	ids := make(map[string]bool)

	for range idCount {
		id := store.SaveReceipt(receipt)
		if id == "" {
			t.Errorf("Expected non-empty ID, got empty string")
		}

		if ids[id] {
			t.Errorf("ID %s was generated more than once", id)
		}

		ids[id] = true
	}
}
