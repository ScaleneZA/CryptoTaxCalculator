package webhandlers

import (
	"encoding/json"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/transactions"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/transactions/ops/calculator"
	"net/http"
)

type CalculateHandler struct {
	B Backends
}

func (c CalculateHandler) Calculate(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var ts []transactions.Transaction
	err := json.NewDecoder(r.Body).Decode(&ts)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// TODO(Make this dynamic from the form)
	fiat := "ZAR"

	yet, err := calculator.Calculate(c.B, fiat, ts)
	if err != nil {
		http.Error(w, "Cannot calculate tax:"+err.Error(), http.StatusBadRequest)
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
