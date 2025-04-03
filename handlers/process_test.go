package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/receipt-processor/models"
	"github.com/receipt-processor/store"
)

func TestProcessHandler(t *testing.T) {
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
		handler := NewProcessHandler(store)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rr.Code)
		}
	})

	t.Run("empty request body", func(t *testing.T) {
		store := store.NewStore()
		handler := NewProcessHandler(store)
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		store := store.NewStore()
		handler := NewProcessHandler(store)
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("invalid json"))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("missing required fields", func(t *testing.T) {
		testCases := []struct {
			name     string
			receipt  models.Receipt
			expected string
		}{
			{"empty retailer", models.Receipt{
				Retailer:     "",
				PurchaseDate: "2023-10-01",
				PurchaseTime: "15:00",
				Items:        []models.Item{{ShortDescription: "Item", Price: "10.00"}},
				Total:        "10.00",
			}, "missing required fields"},
			{"empty date", models.Receipt{
				Retailer:     "Test",
				PurchaseDate: "",
				PurchaseTime: "15:00",
				Items:        []models.Item{{ShortDescription: "Item", Price: "10.00"}},
				Total:        "10.00",
			}, "missing required fields"},
			{"no items", models.Receipt{
				Retailer:     "Test",
				PurchaseDate: "2023-10-01",
				PurchaseTime: "15:00",
				Items:        []models.Item{},
				Total:        "10.00",
			}, "missing required fields"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				store := store.NewStore()
				handler := NewProcessHandler(store)
				body, _ := json.Marshal(tc.receipt)
				req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
				rr := httptest.NewRecorder()

				handler.ServeHTTP(rr, req)

				if rr.Code != http.StatusBadRequest {
					t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
				}
			})
		}
	})

	t.Run("invalid retailer format", func(t *testing.T) {
		invalidReceipt := validReceipt
		invalidReceipt.Retailer = "Invalid@Retailer"
		body, _ := json.Marshal(invalidReceipt)

		store := store.NewStore()
		handler := NewProcessHandler(store)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("invalid date format", func(t *testing.T) {
		invalidReceipt := validReceipt
		invalidReceipt.PurchaseDate = "2023/10/01"
		body, _ := json.Marshal(invalidReceipt)

		store := store.NewStore()
		handler := NewProcessHandler(store)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("invalid time format", func(t *testing.T) {
		invalidReceipt := validReceipt
		invalidReceipt.PurchaseTime = "3:00 PM"
		body, _ := json.Marshal(invalidReceipt)

		store := store.NewStore()
		handler := NewProcessHandler(store)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("invalid total format", func(t *testing.T) {
		testCases := []struct {
			name  string
			total string
		}{
			{"no decimal point", "1000"},
			{"too many decimal places", "10.000"},
			{"non-numeric", "ten dollars"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				invalidReceipt := validReceipt
				invalidReceipt.Total = tc.total
				body, _ := json.Marshal(invalidReceipt)

				store := store.NewStore()
				handler := NewProcessHandler(store)
				req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
				rr := httptest.NewRecorder()

				handler.ServeHTTP(rr, req)

				if rr.Code != http.StatusBadRequest {
					t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
				}
			})
		}
	})

	t.Run("invalid item description", func(t *testing.T) {
		invalidReceipt := validReceipt
		invalidReceipt.Items[0].ShortDescription = "Invalid@Item"
		body, _ := json.Marshal(invalidReceipt)

		store := store.NewStore()
		handler := NewProcessHandler(store)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("invalid item price", func(t *testing.T) {
		invalidReceipt := validReceipt
		invalidReceipt.Items[0].Price = "10.0" // Missing a decimal place
		body, _ := json.Marshal(invalidReceipt)

		store := store.NewStore()
		handler := NewProcessHandler(store)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("successful receipt processing", func(t *testing.T) {
		store := store.NewStore()
		handler := NewProcessHandler(store)
		body, _ := json.Marshal(validReceipt)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		var response models.ReceiptID
		if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}
	})
}

func TestValidateMoneyFormat(t *testing.T) {
	testCases := []struct {
		name    string
		amount  string
		isValid bool
	}{
		{"valid amount", "10.00", true},
		{"no decimal", "10", false},
		{"too many decimals", "10.000", false},
		{"non-numeric", "abc.def", false},
		{"single decimal", "10.0", false},
		{"empty string", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateMoneyFormat(tc.amount)
			if tc.isValid && err != nil {
				t.Errorf("expected valid, got error: %v", err)
			}
			if !tc.isValid && err == nil {
				t.Error("expected invalid, got no error")
			}
		})
	}
}
