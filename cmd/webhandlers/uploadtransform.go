package webhandlers

import (
	"net/http"
	"os"

	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/filetransformer"
)

func UploadTransform(w http.ResponseWriter, r *http.Request) {

}

func transform(filename string) {
	pwd, _ := os.Getwd()

	_, err := filetransformer.Transform(pwd+"/cmd/filetransformer/testData/LUNO_XBT.csv", filetransformer.TransformTypeLuno)
	if err != nil {
		panic(err)
	}
}
