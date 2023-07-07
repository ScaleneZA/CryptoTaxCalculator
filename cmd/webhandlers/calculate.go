package webhandlers

import (
	"encoding/json"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/calculator"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/sharedtypes"
	"net/http"
)

func Calculate(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var ts []sharedtypes.Transaction
	err := json.NewDecoder(r.Body).Decode(&ts)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	yet, err := calculator.Calculate(ts)
	if err != nil {
		http.Error(w, "Cannot calculate tax", http.StatusBadRequest)
		return
	}

	responseJSON, err := json.Marshal(yet)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the response content type
	w.Header().Set("Content-Type", "application/json")

	// Send the response back to the client
	w.Write(responseJSON)
}
