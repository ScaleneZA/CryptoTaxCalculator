package webhandlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/filetransformer"
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

	fmt.Fprintf(w, "Successfully Uploaded File\n")

	ts, err := filetransformer.Transform(file, filetransformer.TransformTypeBasic)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error: Failed to tranform", http.StatusInternalServerError)
	}

	for _, t := range ts {
		fmt.Fprintf(w, "Amount:%f\n", t.Amount)
	}
}

// func transform(filename string) {
// 	pwd, _ := os.Getwd()

// 	_, err := filetransformer.Transform(pwd+"/cmd/filetransformer/testData/LUNO_XBT.csv", filetransformer.TransformTypeLuno)
// 	if err != nil {
// 		panic(err)
// 	}
// }
