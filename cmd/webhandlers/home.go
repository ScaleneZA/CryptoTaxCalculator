package webhandlers

import (
	"github.com/ScaleneZA/CryptoTaxCalculator/cmd/sharedtypes"
	"html/template"
	"log"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/home.html", "web/templates/base.html")
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Define any data you want to pass to the template
	data := struct {
		Title          string
		TransformTypes []sharedtypes.TransformType
	}{
		Title:          "Tax Calculator",
		TransformTypes: sharedtypes.SelectableTransformTypes(),
	}

	// Execute the template with the data
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
