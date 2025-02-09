package main

import (
	"log"
	"net/http"
	"receipt-processor/handlers"
)

func main() {
	// setup endpoint for process receipt - POST method
	http.HandleFunc("/receipts/process", handlers.ProcessReceiptHandler)

	// setup endpoint for get points - GET method
	http.HandleFunc("/receipts/", handlers.GetPointsHandler)

	log.Println("Server is running on port 8080...")
	// start server and check error
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
