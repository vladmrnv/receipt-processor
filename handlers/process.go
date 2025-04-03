package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/receipt-processor/models"
	"github.com/receipt-processor/store"
)

type ProcessHandler struct {
	store *store.Store
}

func NewProcessHandler(s *store.Store) *ProcessHandler {
	return &ProcessHandler{store: s}
}

func (h *ProcessHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	receipt, err := decodeAndValidateReceipt(r)
	if err != nil {
		respondWithError(w, "The receipt is invalid.", http.StatusBadRequest)
		return
	}

	id := h.store.SaveReceipt(receipt)
	respondWithID(w, id)
}

func decodeAndValidateReceipt(r *http.Request) (models.Receipt, error) {
	var receipt models.Receipt
	if r.Body == nil {
		return receipt, ErrEmptyBody
	}

	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		return receipt, err
	}

	if err := validateReceipt(receipt); err != nil {
		return receipt, err
	}

	return receipt, nil
}

func validateReceipt(receipt models.Receipt) error {

	if receipt.Retailer == "" || receipt.PurchaseDate == "" ||
		receipt.PurchaseTime == "" || len(receipt.Items) == 0 ||
		receipt.Total == "" {
		return ErrMissingRequiredFields
	}

	for _, r := range receipt.Retailer {
		if !(r >= 'a' && r <= 'z') && !(r >= 'A' && r <= 'Z') &&
			!(r >= '0' && r <= '9') && r != ' ' && r != '-' && r != '&' {
			return ErrInvalidRetailer
		}
	}

	if _, err := time.Parse("2006-01-02", receipt.PurchaseDate); err != nil {
		return ErrInvalidDate
	}

	if _, err := time.Parse("15:04", receipt.PurchaseTime); err != nil {
		return ErrInvalidTime
	}

	if err := validateMoneyFormat(receipt.Total); err != nil {
		return ErrInvalidTotal
	}

	for _, item := range receipt.Items {
		if strings.TrimSpace(item.ShortDescription) == "" {
			return ErrInvalidItemDescription
		}

		for _, r := range item.ShortDescription {
			if !(r >= 'a' && r <= 'z') && !(r >= 'A' && r <= 'Z') &&
				!(r >= '0' && r <= '9') && r != ' ' && r != '-' {
				return ErrInvalidItemDescription
			}
		}

		if err := validateMoneyFormat(item.Price); err != nil {
			return ErrInvalidItemPrice
		}
	}

	return nil
}

func validateMoneyFormat(amount string) error {
	parts := strings.Split(amount, ".")
	if len(parts) != 2 {
		return errors.New("invalid money format")
	}

	if _, err := strconv.Atoi(parts[0]); err != nil {
		return errors.New("invalid dollars amount")
	}

	if len(parts[1]) != 2 {
		return errors.New("cents must be 2 digits")
	}
	if _, err := strconv.Atoi(parts[1]); err != nil {
		return errors.New("invalid cents amount")
	}

	return nil
}

func respondWithError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func respondWithID(w http.ResponseWriter, id string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.ReceiptID{ID: id})
}

var (
	ErrEmptyBody              = errors.New("empty request body")
	ErrMissingRequiredFields  = errors.New("missing required fields")
	ErrInvalidRetailer        = errors.New("invalid retailer format")
	ErrInvalidDate            = errors.New("invalid purchase date")
	ErrInvalidTime            = errors.New("invalid purchase time")
	ErrInvalidTotal           = errors.New("invalid total amount format")
	ErrInvalidItemDescription = errors.New("invalid item description format")
	ErrInvalidItemPrice       = errors.New("invalid item price format")
)
