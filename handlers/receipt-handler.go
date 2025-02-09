package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"receipt-processor/models"
	"receipt-processor/store"
	"strings"

	"github.com/google/uuid"
)

// ProcessReceiptHandler
func ProcessReceiptHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Only POST Method is Allowed", http.StatusMethodNotAllowed)
		return
	}

	var receipt models.Receipt
	decoder_err := json.NewDecoder(req.Body).Decode(&receipt)
	if decoder_err != nil {
		log.Println("Decoding error:", decoder_err)
		http.Error(w, "Please verify input", http.StatusBadRequest)
		return
	}

	validate_err := store.ValidateReceipt(&receipt)
	if validate_err != nil {
		log.Println("Validation error:", validate_err)
		http.Error(w, "Please verify input", http.StatusBadRequest)
		return
	}

	receipt_id := uuid.New().String()

	total_points := store.CalculatePoints(&receipt)
	store.SaveReceipt(receipt_id, total_points)

	result := map[string]string{"id": receipt_id}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GetPointsHandler
func GetPointsHandler(w http.ResponseWriter, req *http.Request) {
	// Debug: Print the raw URL path
	log.Printf("Received request with path: %s", req.URL.Path)

	if req.Method != http.MethodGet {
		http.Error(w, "Only GET Method is Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from path more reliably
	parts := strings.Split(req.URL.Path, "/")
	log.Printf("Path parts: %v", parts) // Debug: Print path parts

	// Path should be like /receipts/{id}/points
	// So parts would be ["", "receipts", "{id}", "points"]
	if len(parts) < 4 || parts[len(parts)-1] != "points" {
		log.Printf("Invalid path structure") // Debug
		http.Error(w, "Invalid request format. Use /receipts/{id}/points", http.StatusBadRequest)
		return
	}

	// Get the ID (it should be the second-to-last part if 'points' is at the end)
	receipt_id := parts[len(parts)-2]
	log.Printf("Extracted receipt ID: %s", receipt_id) // Debug

	points_value, found := store.GetPoints(receipt_id)
	if !found {
		http.Error(w, "Cannot find receipt with this ID", http.StatusNotFound)
		return
	}

	// prepare and send response
	result := map[string]int{"points": points_value}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
