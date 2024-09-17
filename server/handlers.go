package server

import (
	"fmt"
	"html/template"
	"net/http"

	api "groupie-tracker/api"
)

// Homepage handler that executes the template.html file.
func HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	data := api.ArtistData()

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	fmt.Println("Endpoint: Main page")
}

// Display artist page when an artist's image is clicked, by matching
// the "ArtistName" value against the names in the Data.A.Name field.
func ArtistPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/artistInfo" {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	value := r.FormValue("ArtistName")

	if value == "" {
		ErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	a := api.CollectData()
	var b api.Data
	found := false

	for i, ele := range a {
		if value == ele.A.Name {
			b = a[i]
			Reformat(b)
			found = true
			break
		}
	}

	if !found {
		ErrorHandler(w, r, http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, b); err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	fmt.Println("Endpoint: " + value + "'s page")
}
