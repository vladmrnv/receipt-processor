package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/receipt-processor/handlers"
	"github.com/receipt-processor/store"
)

func main() {
	receiptStore := store.NewStore()

	processHandler := handlers.NewProcessHandler(receiptStore)
	pointsHandler := handlers.NewPointsHandler(receiptStore)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		switch {
		case path == "/receipts/process":
			processHandler.ServeHTTP(w, r)
		case strings.HasPrefix(path, "/receipts/") && strings.HasSuffix(path, "/points"):

			idPart := strings.TrimPrefix(path, "/receipts/")
			id := strings.TrimSuffix(idPart, "/points")

			ctx := context.WithValue(r.Context(), "receipt_id", id)
			pointsHandler.ServeHTTP(w, r.WithContext(ctx))
		default:
			http.NotFound(w, r)
		}
	})

	port := 8080
	fmt.Printf("Server starting on port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
