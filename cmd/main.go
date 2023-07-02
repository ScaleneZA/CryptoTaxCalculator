package main

import (
	"fmt"
	"net/http"

	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/webhandlers"
)

func main() {
	http.HandleFunc("/", webhandlers.Home)
	http.HandleFunc("/ajax/upload_transform", webhandlers.UploadTransform)
	http.HandleFunc("/ajax/calculate", webhandlers.Calculate)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
