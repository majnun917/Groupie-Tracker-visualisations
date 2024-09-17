package server

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	api "groupie-tracker/api"
)

type Text struct {
	ErrorNum int
	ErrorMes string
}

// Custom handler to prevent directory listing
func NoDirListing(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			ErrorHandler(w, r, http.StatusNotFound)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// Reformats information on artist page
func Reformat(b api.Data) api.Data {
	for i, date := range b.D.Dates {
		b.D.Dates[i] = strings.Replace(date, "*", "", -1)
	}

	for j, locus := range b.L.Locations {
		locus = strings.Replace(locus, "_", " ", -1)
		b.L.Locations[j] = strings.Replace(locus, "-", ", ", -1)
	}

	for key, value := range b.R.DatesLocations {
		delete(b.R.DatesLocations, key)
		key := strings.Replace(key, "_", " ", -1)
		key = strings.Replace(key, "-", ", ", -1)
		b.R.DatesLocations[key] = value
	}

	return b
}

// Handles error messages
func ErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)

	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		Error500(w, r, http.StatusInternalServerError)
		return
	}

	var em string
	switch status {
	case http.StatusInternalServerError:
		em = "HTTP status 500: Internal Server Error"
	case http.StatusNotFound:
		em = "HTTP status 404: Page Not Found"
	case http.StatusBadRequest:
		em = "HTTP status 400: Bad Request!\n\nPlease select artist from the Home Page"
	default:
		em = "HTTP status " + strconv.Itoa(status) + ": Something went wrong"
	}

	p := Text{
		ErrorNum: status,
		ErrorMes: em,
	}

	if err := tmpl.Execute(w, p); err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
}

// Handles 500 HTTP status response when error.html is absent
func Error500(w http.ResponseWriter, r *http.Request, status int) {
	fmt.Fprintf(w, `<html>
	<link rel="stylesheet"href="/static/css/styles.css">
	<head><meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Groupie-Tracker</title></head>
	<body><div class="wrapper"><div class="error">
	<h1>%d</h1><pre>We are sorry, but something went wrong.</pre>
	</div><a href="/">Go To Home Page</a></div></body>	
	</html>`, http.StatusInternalServerError)
}
