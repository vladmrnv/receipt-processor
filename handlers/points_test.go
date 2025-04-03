package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/receipt-processor/models"
	"github.com/receipt-processor/processor"
	"github.com/receipt-processor/store"
)

func TestPointsHandler(t *testing.T) {
	validReceipt := models.Receipt{
		Retailer:     "Test Retailer",
		PurchaseDate: "2023-10-01",
		PurchaseTime: "15:00",
		Items: []models.Item{
			{ShortDescription: "Item 1", Price: "10.00"},
		},
		Total: "10.00",
	}

	t.Run("invalid HTTP method", func(t *testing.T) {
		store := store.NewStore()
		handler := NewPointsHandler(store)
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
		}
	})

	t.Run("missing receipt ID", func(t *testing.T) {
		store := store.NewStore()
		handler := NewPointsHandler(store)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("invalid receipt ID type", func(t *testing.T) {
		store := store.NewStore()
		handler := NewPointsHandler(store)
		ctx := context.WithValue(context.Background(), "receipt_id", 123)
		req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("receipt not found", func(t *testing.T) {
		store := store.NewStore()
		handler := NewPointsHandler(store)
		ctx := context.WithValue(context.Background(), "receipt_id", "nonexistent")
		req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
		}

		var response map[string]string
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if response["error"] != "No receipt found for that ID." {
			t.Errorf("unexpected error message: %s", response["error"])
		}
	})

	t.Run("successful points calculation", func(t *testing.T) {
		store := store.NewStore()
		id := store.SaveReceipt(validReceipt)
		expectedPoints := processor.CalculatePoints(validReceipt)

		handler := NewPointsHandler(store)
		ctx := context.WithValue(context.Background(), "receipt_id", id)
		req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
		}

		var response models.Points
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if response.Points != expectedPoints {
			t.Errorf("expected points %d, got %d", expectedPoints, response.Points)
		}
	})
}
