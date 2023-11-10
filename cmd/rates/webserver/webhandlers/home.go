package webhandlers

import (
	"encoding/json"
	"fmt"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/rates/conversionrate/ops/marketvalue"
	"html/template"
	"log"
	"net/http"
	"time"
)

type HomeHandler struct {
	B Backends
}

func (h HomeHandler) Handle(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("cmd/rates/webserver/templates/home.html", "cmd/rates/webserver/templates/base.html")
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Define any data you want to pass to the template
	data := struct {
		Title string
	}{
		Title: "Rates",
	}

	// Execute the template with the data
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h HomeHandler) ClosestForTime(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the incoming JSON data
	var data struct {
		Date string `json:"date"`
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	t, err := time.Parse(time.DateOnly, data.Date)
	if err != nil {
		http.Error(w, "Failed to parse date", http.StatusBadRequest)
		return
	}

	mps, err := marketvalue.ClosestMarketPairsAtPoint(h.B, t.Unix())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error finding closest for timestamp: %d", t.Unix()), http.StatusInternalServerError)
		return
	}
	responseData := struct {
		Message string                      `json:"message"`
		Data    []conversionrate.MarketPair `json:"data"`
	}{
		Message: "success",
		Data:    mps,
	}

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
