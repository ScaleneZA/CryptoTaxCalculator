package webhandlers

import (
	"encoding/json"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/transactions"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/transactions/ops/calculator"
	"github.com/luno/jettison/log"
	"html/template"
	"net/http"
)

type CalculateHandler struct {
	B Backends
}

func subtract(a, b int) int {
	return a - b
}

func (c CalculateHandler) Calculate(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("cmd/transactions/webserver/templates/taxpacks.html")
	if err != nil {
		log.Error(r.Context(), err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl.Funcs(template.FuncMap{
		"subtract": subtract,
	})

	var ts []transactions.Transaction
	err = json.NewDecoder(r.Body).Decode(&ts)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// TODO(Make this dynamic from the form)
	fiat := "ZAR"

	yet, err := calculator.Calculate(c.B, fiat, ts)
	if err != nil {
		http.Error(w, "Cannot calculate tax: "+err.Error(), http.StatusBadRequest)
		return
	}

	data := struct {
		Title         string
		YearEndTotals []calculator.YearEndTotal
	}{
		Title:         "Tax Calculator - Tax packs",
		YearEndTotals: yet,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
