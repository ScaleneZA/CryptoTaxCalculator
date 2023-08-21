package webhandlers

import (
	"encoding/json"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/transactions"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/transactions/db/calculator"
	"net/http"
)

type OverrideHandler struct {
	B Backends
}

type JsonResponse struct {
	Message string `json:"message"`
}

type JsonRequest struct {
	UID           string                       `json:"uid""`
	OverridedType transactions.TransactionType `json:"overrided_type"`
}

func (o OverrideHandler) Override(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the incoming JSON data
	var data JsonRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	_, err = calculator.Upsert(o.B.DB(), data.UID, data.OverridedType)
	if err != nil {
		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
	}

	responseData := JsonResponse{Message: "success"}

	// Set the response Content-Type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode and send the response JSON
	encoder := json.NewEncoder(w)
	err = encoder.Encode(responseData)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}
