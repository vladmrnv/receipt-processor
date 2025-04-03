package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/receipt-processor/models"
	"github.com/receipt-processor/processor"
	"github.com/receipt-processor/store"
)

type PointsHandler struct {
	Store *store.Store
}

func NewPointsHandler(s *store.Store) *PointsHandler {
	return &PointsHandler{Store: s}
}

func (h *PointsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, ok := r.Context().Value("receipt_id").(string)
	if !ok {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	receipt, err := h.Store.GetReceipt(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "No receipt found for that ID.",
		})
		return
	}

	points := processor.CalculatePoints(receipt)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Points{Points: points})
}
