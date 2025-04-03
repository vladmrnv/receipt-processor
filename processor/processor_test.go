package processor

import (
	"testing"

	"github.com/receipt-processor/models"
)

func TestCalculatePoints(t *testing.T) {
	tests := []struct {
		name     string
		receipt  models.Receipt
		expected int
	}{
		{
			name: "Target Receipt Example",
			receipt: models.Receipt{
				Retailer:     "Target",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items: []models.Item{
					{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
					{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
					{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
					{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
					{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
				},
				Total: "35.35",
			},
			expected: 28,
		},
		{
			name: "Corner Market Receipt Example",
			receipt: models.Receipt{
				Retailer:     "M&M Corner Market",
				PurchaseDate: "2022-03-20",
				PurchaseTime: "14:33",
				Items: []models.Item{
					{ShortDescription: "Gatorade", Price: "2.25"},
					{ShortDescription: "Gatorade", Price: "2.25"},
					{ShortDescription: "Gatorade", Price: "2.25"},
					{ShortDescription: "Gatorade", Price: "2.25"},
				},
				Total: "9.00",
			},
			expected: 109,
		},
		{
			name: "Round Dollar Amount",
			receipt: models.Receipt{
				Retailer:     "ABC",
				PurchaseDate: "2022-01-02",
				PurchaseTime: "12:00",
				Items: []models.Item{
					{ShortDescription: "Item", Price: "10.00"},
				},
				Total: "10.00",
			},

			expected: 78,
		},
		{
			name: "Multiple of Quarter",
			receipt: models.Receipt{
				Retailer:     "XYZ",
				PurchaseDate: "2022-01-03",
				PurchaseTime: "15:30",
				Items: []models.Item{
					{ShortDescription: "Item 1", Price: "5.25"},
					{ShortDescription: "Item 2", Price: "5.25"},
				},
				Total: "10.50",
			},
			expected: 53,
		},
		{
			name: "Alphanumeric Characters in Retailer",
			receipt: models.Receipt{
				Retailer:     "Shop123",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "12:00",
				Items: []models.Item{
					{ShortDescription: "Item", Price: "1.00"},
				},
				Total: "1.00",
			},

			expected: 88,
		},
		{
			name: "Item Description Multiple of 3",
			receipt: models.Receipt{
				Retailer:     "ABC",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "12:00",
				Items: []models.Item{
					{ShortDescription: "Gum", Price: "2.25"},
					{ShortDescription: "Coffee", Price: "3.50"},
					{ShortDescription: "Tea Bag", Price: "1.25"},
				},
				Total: "7.00",
			},
			expected: 91,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			points := CalculatePoints(tc.receipt)
			if points != tc.expected {
				t.Errorf("Expected %d points, got %d", tc.expected, points)
			}
		})
	}
}

func TestCountAlphanumeric(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"ABC123", 6},
		{"Target", 6},
		{"M&M Corner Market", 14},
		{"Store-123", 8},
		{"", 0},
		{"$%@#@", 0},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			count := countAlphanumeric(tc.input)
			if count != tc.expected {
				t.Errorf("Expected %d alphanumeric characters, got %d", tc.expected, count)
			}
		})
	}
}
