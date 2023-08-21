package webhandlers

import (
	"encoding/json"
	"net/http"
)

type OverrideHandler struct {
	B Backends
}

type JsonResponse struct {
	Message string `json:"message"`
}

func (o OverrideHandler) Override(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the incoming JSON data
	var requestData map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	responseData := JsonResponse{Message: "Received your POST request"}

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
