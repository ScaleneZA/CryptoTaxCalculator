package webhandlers

import (
	"encoding/json"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/transactions"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/transactions/ops/calculator"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/transactions/transactions/ops/filetransformer"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
)

type UploadTransformer struct {
	B Backends
}

func (ut UploadTransformer) UploadTransform(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	fs := make(map[transactions.TransformType][]io.Reader)
	for _, typ := range transactions.SelectableTransformTypes() {
		files, ok := r.MultipartForm.File["files-"+strconv.Itoa(int(typ))]
		if !ok {
			continue
		}
		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				log.Println(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			defer file.Close()

			log.Printf("Uploaded File: %+v\n", fileHeader.Filename)
			log.Printf("File Size: %+v\n", fileHeader.Size)
			log.Printf("MIME Header: %+v\n", fileHeader.Header)

			fs[typ] = append(fs[typ], file)
		}
	}

	ts, err := filetransformer.TransformAll(fs)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error: Failed to tranform", http.StatusInternalServerError)
	}

	ts, err = calculator.PopulateOverriddenTypes(ut.B, ts)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error: Failed to populate override types", http.StatusInternalServerError)
	}

	jsonData, err := json.Marshal(ts)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("cmd/transactions/webserver/templates/overrides.html", "cmd/transactions/webserver/templates/base.html")
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title                    string
		OverrideTransactionTypes []transactions.TransactionType
		Transactions             []transactions.Transaction
		TransactionsJSON         string
	}{
		Title:                    "Overrides",
		OverrideTransactionTypes: transactions.ValidTransactionTypes(),
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
