package webhandlers

import (
	"html/template"
	"log"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/home.html")
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Define any data you want to pass to the template
	data := struct {
		Title string
	}{
		Title: "My Web App",
	}

	// Execute the template with the data
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
