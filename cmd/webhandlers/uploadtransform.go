package webhandlers

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/filetransformer"
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/sharedtypes"
)

func UploadTransform(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var fs []io.Reader
	files := r.MultipartForm.File["files"]
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

		fs = append(fs, file)
	}

	typInt, err := strconv.Atoi(r.FormValue("type"))
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error: Failed to parse type", http.StatusInternalServerError)
	}

	ts, err := filetransformer.Transform(fs, filetransformer.TransformType(typInt))
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
