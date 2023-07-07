package webhandlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/filetransformer"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/sharedtypes"
)

func UploadTransform(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	defer file.Close()
	log.Printf("Uploaded File: %+v\n", handler.Filename)
	log.Printf("File Size: %+v\n", handler.Size)
	log.Printf("MIME Header: %+v\n", handler.Header)

	typInt, err := strconv.Atoi(r.FormValue("type"))
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error: Failed to parse type", http.StatusInternalServerError)
	}

	ts, err := filetransformer.Transform(file, filetransformer.TransformType(typInt))
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error: Failed to tranform", http.StatusInternalServerError)
	}

	jsonData, err := json.Marshal(ts)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("web/templates/overrides.html", "web/templates/base.html")
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title                    string
		OverrideTransactionTypes []sharedtypes.TransactionType
		Transactions             []sharedtypes.Transaction
		TransactionsJSON         string
	}{
		Title:                    "Overrides",
		OverrideTransactionTypes: sharedtypes.ValidTransactionTypes(),
		Transactions:             ts,
		TransactionsJSON:         string(jsonData),
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
